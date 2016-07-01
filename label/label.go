package label

// DataLabel is the label data stored in the database
type DataLabel struct {
	ID   int64
	Name string
}

// Label includes DataLabel and any computed values
type Label struct {
	DataLabel
}
