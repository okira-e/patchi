package types

type TableDiff struct {
	TableName string
	// "created" | "updated" | "deleted"
	diffType string
}
