package dbrutils

import (
	"reflect"

	"github.com/gocraft/dbr"
)

// EqExpr creates an equality comparison expression for columns from two table.
func EqExpr(table1, column1, table2, column2 string) dbr.Builder {
	return dbr.Expr("? = ?", I(table1, column1), I(table2, column2))
}

// Eq represents an equality comparison.
// This method enhance dbr.Eq for allowing passing column in any type.
func Eq(column interface{}, value interface{}) dbr.Builder {
	switch c := column.(type) {
	case string:
		return dbr.Eq(c, value)
	case dbr.Builder:
		break
	default:
		if r := reflect.ValueOf(column); r.Kind() == reflect.String {
			return dbr.Eq(r.String(), value)
		}
	}

	if value == nil {
		return dbr.Expr("? IS NULL", column)
	}
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Slice {
		if size := v.Len(); size == 0 {
			return falsy()
		} else if size > 1 {
			return dbr.Expr("? IN ?", column, value)
		}
		value = v.Index(0).Interface()
	}
	return dbr.Expr("? = ?", column, value)
}

func falsy() dbr.Builder {
	return dbr.BuildFunc(func(d dbr.Dialect, buf dbr.Buffer) error {
		buf.WriteString(d.EncodeBool(false))
		return nil
	})
}

// MatchInNaturalLang creates a match comparison against full text index column in natural language mode.
func MatchInNaturalLang(column interface{}, value interface{}) dbr.Builder {
	return dbr.Expr("MATCH (?) AGAINST (? IN NATURAL LANGUAGE MODE)", column, value)
}
