package label

const (
	// Key is the format by which label keys are made
	Key = "label_%d"
	// IDByName is the index to look up label ids by name
	IDByName = "labelName_%s"

	// CountKey is the key of the auto-incrementer
	CountKey = "labelCounter"
)

// DataLabel is the label data stored in the database
type DataLabel struct {
	ID   int64
	Name string
}

// Label includes DataLabel and any computed values
type Label struct {
	DataLabel
}

// Populate fills in any computed values
func (l *Label) Populate() {
}
