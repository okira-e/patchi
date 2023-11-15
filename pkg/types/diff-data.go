package types

type TableDiff struct {
	TableName string
	// DiffType ("created" | "updated" | "deleted") represents the type of change that has occurred to the table.
	DiffType string
}
