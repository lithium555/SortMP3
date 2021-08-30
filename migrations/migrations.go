package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func main() {
	if err := upMigrations(); err != nil {
		log.Fatalf("Run(). Error: '%v'", err)
	}
}

//go:embed files/*.sql
var fs embed.FS

func upMigrations() error {
	pgConString := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	d, err := iofs.New(fs, "files")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, pgConString)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil {
		return err
	}

	return nil
}
