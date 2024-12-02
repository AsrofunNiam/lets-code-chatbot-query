package domain

type SchemaTables []SchemaTable
type SchemaTable struct {
	TableName  string
	ColumnName string
	DataType   string
}
