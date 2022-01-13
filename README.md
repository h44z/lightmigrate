# LightMigrate - a lightweight database migration library

[![codecov](https://codecov.io/gh/h44z/lightmigrate/branch/master/graph/badge.svg?token=MMAOVBLL2U)](https://codecov.io/gh/h44z/lightmigrate)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![GoDoc](https://pkg.go.dev/badge/github.com/h44z/lightmigrate)](https://pkg.go.dev/github.com/h44z/lightmigrate)
![GitHub last commit](https://img.shields.io/github/last-commit/h44z/lightmigrate)
[![Go Report Card](https://goreportcard.com/badge/github.com/h44z/lightmigrate)](https://goreportcard.com/report/github.com/h44z/lightmigrate)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/h44z/lightmigrate)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/h44z/lightmigrate)
[![GitHub Release](https://img.shields.io/github/release/h44z/lightmigrate.svg)](https://github.com/h44z/lightmigrate/releases)

This library is heavily inspired by [golang-migrate](https://github.com/golang-migrate/migrate).

It currently lacks support for many database drivers and the CLI feature, this is still WIP. 

But it is completely restructured to minimize the dependency footprint.

## Currently Supported databases
 - [MongoDB](https://github.com/h44z/lightmigrate-mongodb) 

## Usage example:

```go
fsys := os.DirFS("/app/data") // Migration File Source Root (/app/data/migration-files)

source, err := NewFsSource(fsys, "migration-files")
if err != nil {
    t.Fatalf("unable to setup source: %v", err)
}
defer source.Close()

driver, err := test.NewMockDriver() // Database Driver (for example lightmigrate_mongo.NewDriver())
if err != nil {
    t.Fatalf("unable to setup driver: %v", err)
}
defer driver.Close()

migrator, err := NewMigrator(source, driver, WithVerboseLogging(true)) // The migrator instance
if err != nil {
    t.Fatalf("unable to setup migrator: %v", err)
}

err = migrator.Migrate(3) // Migrate to schema version 3
if err != nil {
    t.Fatalf("migration error: %v", err)
}
```

