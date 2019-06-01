// update-orders-sheet updates the orders spreadsheet in Google Sheets to
// reflect a new or revised order.
//
// usage: update-orders-sheet orderID

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"syscall"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/sheets/v4"

	"scholacantorum.org/orders/config"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

const sheet = "1o34FdjNIrzMUPuZzR7d1bfGWtzleUFqJl7BKDvuRkT0"

func main() {
	var (
		logfile  *os.File
		orderID  model.OrderID
		order    *model.Order
		tx       db.Tx
		lockfh   *os.File
		conf     *jwt.Config
		sheet    string
		svc      *sheets.Service
		ss       *sheets.Spreadsheet
		sheetnum int64
		vr       *sheets.BatchGetValuesResponse
		requests []*sheets.Request
		lines    = map[model.ProductID]*model.OrderLine{}
		client   *http.Client
		err      error
	)
	// Initialize the logger.  Since we expect it to exist, this will also
	// confirm that we're in the data directory.
	if logfile, err = os.OpenFile("server.log", os.O_APPEND|os.O_WRONLY, 0600); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	log.SetOutput(logfile)
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetPrefix("update-orders-sheet")
	// Log any panics.
	defer func() {
		if panicked := recover(); panicked != nil {
			log.Printf("PANIC: %v", panicked)
			fmt.Fprint(logfile, string(debug.Stack()))
			os.Exit(1)
		}
	}()
	// Get the order ID from the command line.
	if len(os.Args) == 2 {
		if id, err := strconv.Atoi(os.Args[1]); err == nil && id > 0 {
			orderID = model.OrderID(id)
		}
	}
	if orderID == 0 {
		log.Fatalf("usage: update-orders-sheet orderID")
	}
	// Open the database and read the order.
	db.Open("orders.db")
	tx = db.Begin()
	if order = tx.FetchOrder(orderID); order == nil {
		log.Fatalf("order %d does not exist", orderID)
	}
	tx.Commit()
	// Obtain a file lock to ensure that we don't have parallel spreadsheet
	// updates.
	if lockfh, err = os.OpenFile(config.Get("sheetLockFile"), os.O_CREATE|os.O_RDWR, 0644); err != nil {
		log.Fatal(err)
	}
	defer lockfh.Close()
	if err = syscall.Flock(int(lockfh.Fd()), syscall.LOCK_EX); err != nil {
		log.Fatal(err)
	}
	// Establish a connection to Google Sheets, as admin@scholacantorum.org.
	conf = &jwt.Config{
		Email:      config.Get("sheetsEmail"),
		PrivateKey: []byte(config.Get("sheetsPrivateKey")),
		Scopes:     []string{"https://www.googleapis.com/auth/spreadsheets"},
		TokenURL:   google.JWTTokenURL,
		Subject:    "admin@scholacantorum.org",
	}
	client = conf.Client(oauth2.NoContext)
	if svc, err = sheets.New(client); err != nil {
		log.Fatal(err)
	}
	sheet = config.Get("sheetID")
	// Get the sheet number for the "Orders" sheet.
	if ss, err = svc.Spreadsheets.Get(sheet).Fields(googleapi.Field("sheets.properties")).Do(); err != nil {
		log.Fatal(err)
	}
	for _, sh := range ss.Sheets {
		if sh.Properties.Title == "Orders" {
			sheetnum = sh.Properties.SheetId
			break
		}
	}
	if sheetnum == 0 {
		err = fmt.Errorf(`no "Orders" sheet found in spreadsheet`)
		log.Fatal(err)
	}
	// Get the order IDs and products from columns A and K.
	if vr, err = svc.Spreadsheets.Values.BatchGet(sheet).Ranges("Orders!A:A", "Orders!K:K").Do(); err != nil {
		log.Fatal(err)
	}
	// If the order isn't valid, we'll act as if it has no lines.
	if order.Flags&model.OrderValid == 0 {
		order.Lines = nil
	}
	// Make a map of lines by product.
	for _, ol := range order.Lines {
		lines[ol.Product.ID] = ol
	}
	// Walk the list of rows in the spreadsheet, looking for ones that
	// belong to this order.  We walk from the bottom up so that deletions
	// don't change row numbers we care about.
	for row := len(vr.ValueRanges[0].Values) - 1; row >= 0; row-- {
		var rowoid int
		var product string
		var line *model.OrderLine
		var ok bool

		if len(vr.ValueRanges[0].Values[row]) < 1 {
			continue
		}
		rowoid, ok = vr.ValueRanges[0].Values[row][0].(int)
		if !ok || rowoid != int(order.ID) {
			continue
		}
		if len(vr.ValueRanges[1].Values[row]) < 1 {
			continue
		}
		product, ok = vr.ValueRanges[1].Values[row][0].(string)
		if ok {
			line, ok = lines[model.ProductID(product)]
		}
		if !ok {
			// This row is for our order, with an unknown product —
			// probably something that was returned.  We want to
			// delete the row.
			requests = append(requests, &sheets.Request{DeleteDimension: &sheets.DeleteDimensionRequest{
				Range: &sheets.DimensionRange{
					Dimension:  "ROWS",
					StartIndex: int64(row),
					EndIndex:   int64(row) + 1,
					SheetId:    sheetnum,
				},
			}})
			continue
		}
		// We found a row for this product.  Update the quantity, unit
		// price, and total on it.
		requests = append(requests, &sheets.Request{UpdateCells: &sheets.UpdateCellsRequest{
			Start: &sheets.GridCoordinate{
				SheetId:     sheetnum,
				RowIndex:    int64(row),
				ColumnIndex: 12, // M, zero based
			},
			Fields: "userEnteredValue",
			Rows: []*sheets.RowData{{Values: []*sheets.CellData{{
				UserEnteredValue: &sheets.ExtendedValue{
					NumberValue: float64(line.Quantity),
				},
			}}}},
		}})
		requests = append(requests, &sheets.Request{UpdateCells: &sheets.UpdateCellsRequest{
			Start: &sheets.GridCoordinate{
				SheetId:     sheetnum,
				RowIndex:    int64(row),
				ColumnIndex: 13, // N, zero based
			},
			Fields: "userEnteredValue",
			Rows: []*sheets.RowData{{Values: []*sheets.CellData{{
				UserEnteredValue: &sheets.ExtendedValue{NumberValue: float64(line.Price / 100)},
			}}}},
		}})
		requests = append(requests, &sheets.Request{UpdateCells: &sheets.UpdateCellsRequest{
			Start: &sheets.GridCoordinate{
				SheetId:     sheetnum,
				RowIndex:    int64(row),
				ColumnIndex: 14, // O, zero based
			},
			Fields: "userEnteredValue",
			Rows: []*sheets.RowData{{Values: []*sheets.CellData{{
				UserEnteredValue: &sheets.ExtendedValue{NumberValue: float64(line.Quantity * line.Price / 100)},
			}}}},
		}})
		delete(lines, line.Product.ID)
	}

	// If there are remaining products that we didn't see, append them to
	// the bottom of the sheet.
	for _, line := range lines {
		requests = append(requests, &sheets.Request{AppendCells: &sheets.AppendCellsRequest{
			SheetId: sheetnum,
			Fields:  "userEnteredValue",
			Rows: []*sheets.RowData{{Values: []*sheets.CellData{{
				UserEnteredValue: &sheets.ExtendedValue{NumberValue: float64(order.ID)},
			}, {
				UserEnteredValue: &sheets.ExtendedValue{StringValue: order.Created.Format("2006-01-02 15:04:05")},
			}, {
				UserEnteredValue: &sheets.ExtendedValue{StringValue: order.Payments[0].Method},
			}, {
				UserEnteredValue: &sheets.ExtendedValue{StringValue: order.Payments[0].Stripe},
			}, {
				UserEnteredValue: &sheets.ExtendedValue{StringValue: order.Name},
			}, {
				UserEnteredValue: &sheets.ExtendedValue{StringValue: order.Email},
			}, {
				UserEnteredValue: &sheets.ExtendedValue{StringValue: order.Address},
			}, {
				UserEnteredValue: &sheets.ExtendedValue{StringValue: order.City},
			}, {
				UserEnteredValue: &sheets.ExtendedValue{StringValue: order.State},
			}, {
				UserEnteredValue: &sheets.ExtendedValue{StringValue: order.Zip},
			}, {
				UserEnteredValue: &sheets.ExtendedValue{StringValue: string(line.Product.ID)},
			}, {
				UserEnteredValue: &sheets.ExtendedValue{StringValue: string(line.Product.ID)},
			}, {
				UserEnteredValue: &sheets.ExtendedValue{NumberValue: float64(line.Quantity)},
			}, {
				UserEnteredValue: &sheets.ExtendedValue{NumberValue: float64(line.Price / 100)},
			}, {
				UserEnteredValue: &sheets.ExtendedValue{NumberValue: float64(line.Quantity * line.Price / 100)},
			}}}},
		}})
	}
	// Execute the batched change.
	if len(requests) != 0 {
		_, err = svc.Spreadsheets.BatchUpdate(sheet, &sheets.BatchUpdateSpreadsheetRequest{Requests: requests}).Do()
		if err != nil {
			log.Fatal(err)
		}
	}
}
