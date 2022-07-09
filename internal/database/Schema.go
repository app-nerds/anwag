package database

type TableStruct struct {
	Name    string
	Columns []TableColumn
}

type TableColumn struct {
	Name     string
	DataType string
}
