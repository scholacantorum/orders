package model

import (
	"database/sql"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func panicOnExecError(_ sql.Result, err error) {
	if err != nil {
		panic(err)
	}
}

func panicOnNoRows(res sql.Result, err error) {
	var rows int64

	if err != nil {
		panic(err)
	}
	if rows, err = res.RowsAffected(); err != nil {
		panic(err)
	}
	if rows == 0 {
		panic("affected no rows")
	}
}
