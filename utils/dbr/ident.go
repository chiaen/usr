package dbrutils

import "github.com/gocraft/dbr"

// I identity helper function for column with table prefix.
func I(table string, col string) dbr.I {
	return dbr.I(table + "." + col)
}
