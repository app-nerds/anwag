package database

var PostgresDatatypes = map[string]string{
	"bigint":                      "int",
	"bit":                         "bool",
	"character varying":           "string",
	"integer":                     "int",
	"numeric":                     "float64",
	"smallint":                    "int",
	"text":                        "string",
	"time without time zone":      "time.Time",
	"timestamp without time zone": "time.Time",
	"uuid":                        "string",
}
