//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	linux   = map[string]string{"GOOS": "linux"}
	goBuild = sh.RunWith(linux, mg.GoCmd(), "build")
)

var Default = Build

func Build() {
	mg.Deps(UpdateOrdersSheet, ResendReceipt, OrdersAPI)
}

func UpdateOrdersSheet() error {
	return sh.RunWith(linux, mg.GoCmd(), "build", "-o", "dist/update-orders-sheet", "./cmd/update-orders-sheet")
}

func ResendReceipt() error {
	return sh.RunWith(linux, mg.GoCmd(), "build", "-o", "dist/resend-receipt", "./cmd/resend-receipt")
}

func OrdersAPI() error {
	if err := sh.RunWith(linux, mg.GoCmd(), "build", "-o", "dist/ofcapi", "."); err != nil {
		return err
	}
	sh.Run("ln", "-f", "dist/ofcapi", "dist/payapi")
	sh.Run("ln", "-f", "dist/ofcapi", "dist/posapi")
	sh.Run("ln", "-f", "dist/ofcapi", "dist/ticket")
	return nil
}

func InstallSandbox() error {
	mg.Deps(Build)
	if err := sh.Run("scp", "dist/update-orders-sheet", "dist/resend-receipt", "schola:bin"); err != nil {
		return err
	}
	if err := sh.Run("scp", "dist/ofcapi", "schola:orders-test.scholacantorum.org"); err != nil {
		return err
	}
	if err := sh.Run("ssh", "schola", "cd orders-test.scholacantorum.org && ln -f ofcapi payapi && ln -f ofcapi posapi && ln -f ofcapi ticket"); err != nil {
		return err
	}
	return nil
}

func InstallProduction() error {
	mg.Deps(Build)
	if err := sh.Run("scp", "dist/update-orders-sheet", "dist/resend-receipt", "schola:bin"); err != nil {
		return err
	}
	if err := sh.Run("scp", "dist/ofcapi", "schola:orders.scholacantorum.org"); err != nil {
		return err
	}
	if err := sh.Run("ssh", "schola", "cd orders.scholacantorum.org && ln -f ofcapi payapi && ln -f ofcapi posapi && ln -f ofcapi ticket"); err != nil {
		return err
	}
	return nil
}
