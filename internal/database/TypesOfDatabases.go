package database

import "fmt"

const (
	MongoDB  string = "MongoDB"
	MySQL    string = "MySQL"
	Postgres string = "Postgres"
)

func GetDSN(databaseType, dbName string) string {
	switch databaseType {
	case MongoDB:
		return "localhost:27017"

	case MySQL:
		return fmt.Sprintf("root:password@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbName)

	case Postgres:
		return fmt.Sprintf("host=localhost user=postgres password=password dbname=%s port=5432", dbName)
	}

	return ""
}
