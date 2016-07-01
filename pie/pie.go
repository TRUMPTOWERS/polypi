package pie

const (
	// Key is the format by which the keys for pies are formed
	Key = "pie_%d"

	// CountKey is the key use to store the id incrementer
	CountKey = "pieCounter"
)

// DataPie is the pie data stored in the database
type DataPie struct {
	ID        int64
	Name      string
	ImageURL  string
	Slices    int
	Price     float64
	Purchases []int64
	Labels    []int64
}

// Pie includes DataPie, and any computed values
type Pie struct {
	DataPie
	RemainingSlices int
}

// Populate fills in computed values
func (dp *Pie) Populate() {
	dp.RemainingSlices = dp.Slices - len(dp.Purchases)
}
