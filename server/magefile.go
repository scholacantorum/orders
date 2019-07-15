// +build mage

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var goBuild = sh.RunCmd(mg.GoCmd(), "build")
var goInstall = sh.RunCmd(mg.GoCmd(), "install")

var Default = Sandbox

func Sandbox() {
	mg.Deps(UpdateOrdersSheet, ResendReceipt, SandboxOrders, SandboxPaymentForms, SandboxAtTheDoor)
}

func Production() {
	mg.Deps(UpdateOrdersSheet, ResendReceipt, ProductionOrders, ProductionPaymentForms, ProductionAtTheDoor)
}

func UpdateOrdersSheet() error {
	return goInstall("./cmd/update-orders-sheet")
}

func ResendReceipt() error {
	return goInstall("./cmd/resend-receipt")
}

func SandboxOrders() error {
	return orders("/home/scholacantorum/orders-test.scholacantorum.org")
}

func ProductionOrders() error {
	return orders("/home/scholacantorum/orders.scholacantorum.org")
}

func orders(path string) error {
	if err := goBuild("-o", path+"/api", "."); err != nil {
		return err
	}
	if err := sh.Run("ln", "-f", path+"/api", path+"/ticket"); err != nil {
		return err
	}
	return nil
}

func SandboxPaymentForms() error {
	if err := publicPaymentForms(
		"/home/scholacantorum/orders-test.scholacantorum.org",
		"/home/scholacantorum/schola6p/data/sandbox/resources",
		"sandbox.yaml",
	); err != nil {
		return err
	}
	return membersPaymentForms("/home/scholacantorum/orders-test.scholacantorum.org")
}

func ProductionPaymentForms() error {
	if err := publicPaymentForms(
		"/home/scholacantorum/orders.scholacantorum.org",
		"/home/scholacantorum/schola6p/data/production/resources",
		"production.yaml",
	); err != nil {
		return err
	}
	return membersPaymentForms("/home/scholacantorum/orders.scholacantorum.org")
}

func publicPaymentForms(path, data, config string) error {
	var (
		matches   []string
		resfile   *os.File
		savedir   string
		err       error
		resources = map[string]map[string][]string{}
	)
	for _, base := range []string{"chunk-vendors", "buy-tickets", "donate"} {
		matches, _ = filepath.Glob(path + "/" + base + ".*")
		for _, match := range matches {
			if err = sh.Rm(match); err != nil {
				return err
			}
		}
		resources[base] = map[string][]string{}
		for _, ext := range []string{"js", "css"} {
			matches, _ = filepath.Glob(fmt.Sprintf("../payment-forms/dist/%s/%s.*.%s", ext, base, ext))
			for _, match := range matches {
				if err = sh.Copy(path+"/"+filepath.Base(match), match); err != nil {
					return err
				}
				if err = sh.Run("gzip", "-k", path+"/"+filepath.Base(match)); err != nil {
					return err
				}
				// This is a hack for Apache.  If someone
				// requests the file with Accept-Encoding:gzip,
				// it will reply with the one with the .gz
				// extension -- but only if the exact filename
				// they requested doesn't exist.  Fortunately,
				// if we double the extension on the original
				// file, everything works as desired.
				if err = sh.Run("mv", path+"/"+filepath.Base(match),
					fmt.Sprintf("%s/%s.%s", path, filepath.Base(match), ext)); err != nil {
					return err
				}
				resources[base][ext] = append(resources[base][ext], filepath.Base(match))
			}
		}
	}
	if err = sh.Run("mkdir", "-p", data); err != nil {
		return err
	}
	for base := range resources {
		if base != "chunk-vendors" {
			for ext := range resources[base] {
				resources[base][ext] = append(resources[base][ext], resources["chunk-vendors"][ext]...)
			}
			if resfile, err = os.Create(data + "/" + base + ".json"); err != nil {
				return err
			}
			json.NewEncoder(resfile).Encode(resources[base])
			resfile.Close()
		}
	}
	savedir, _ = os.Getwd()
	if err = os.Chdir("/home/scholacantorum/schola6p"); err != nil {
		return err
	}
	if err = sh.Run("hugo", "--config", config); err != nil {
		return err
	}
	if err = os.Chdir(savedir); err != nil {
		return err
	}
	return nil
}

func membersPaymentForms(path string) error {
	var (
		matches []string
		err     error
	)
	for _, base := range []string{"/js/chunk-vendors.*", "/css/chunk-vendors.*", "/js/recordings.*", "/css/recordings.*"} {
		matches, _ = filepath.Glob(path + base)
		for _, match := range matches {
			if err = sh.Rm(match); err != nil {
				return err
			}
		}
	}
	for _, base := range []string{"chunk-vendors", "recordings"} {
		for _, ext := range []string{"js", "css"} {
			matches, _ = filepath.Glob(fmt.Sprintf("../payment-forms/dist/%s/%s.*.%s", ext, base, ext))
			for _, match := range matches {
				if err = sh.Run("mkdir", "-p", path+"/"+ext); err != nil {
					return err
				}
				if err = sh.Copy(path+"/"+ext+"/"+filepath.Base(match), match); err != nil {
					return err
				}
				if err = sh.Run("gzip", "-k", path+"/"+ext+"/"+filepath.Base(match)); err != nil {
					return err
				}
				// Same hack for Apache.
				if err = sh.Run("mv", path+"/"+ext+"/"+filepath.Base(match),
					fmt.Sprintf("%s/%s/%s.%s", path, ext, filepath.Base(match), ext)); err != nil {
					return err
				}
			}
		}
	}
	return sh.Copy(path+"/recordings.html", "../payment-forms/dist/recordings.html")
}

func SandboxAtTheDoor() error {
	return atTheDoor("/home/scholacantorum/orders-test.scholacantorum.org/door")
}

func ProductionAtTheDoor() error {
	return atTheDoor("/home/scholacantorum/orders.scholacantorum.org/door")
}

func atTheDoor(path string) error {
	if err := sh.Rm(path); err != nil {
		return err
	}
	return sh.Run("cp", "-rp", "../door/dist", path)
}

func ResetSandbox() error {
	return reset("/home/scholacantorum/orders-test.scholacantorum.org")
}

func ResetProduction() error {
	return reset("/home/scholacantorum/orders.scholacantorum.org")
}

func reset(path string) error {
	if err := sh.Run("mkdir", "-p", path+"/data"); err != nil {
		return err
	}
	if err := sh.Run("chmod", "700", path+"/data"); err != nil {
		return err
	}
	if err := sh.Rm(path + "/data/orders.db"); err != nil {
		return err
	}
	if err := sh.Rm(path + "/data/server.log"); err != nil {
		return err
	}
	if err := sh.Run("sqlite3", path+"/data/orders.db", ".read db/schema.sql"); err != nil {
		return err
	}
	if err := sh.Run("chmod", "600", path+"/data/orders.db"); err != nil {
		return err
	}
	if err := sh.Run("sqlite3", path+"/data/orders.db", ".read db/seed.sql"); err != nil {
		return err
	}
	if err := sh.Run("touch", path+"/data/server.log"); err != nil {
		return err
	}
	if err := sh.Run("chmod", "600", path+"/data/server.log"); err != nil {
		return err
	}
	sh.Run("chmod", "600", path+"/data/config.json")
	return nil
}
