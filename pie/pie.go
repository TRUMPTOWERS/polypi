package pie

import (
	"github.com/TRUMPTOWERS/polypi/label"
	"github.com/TRUMPTOWERS/polypi/purchase"
)

// DataPie is the pie data stored in the database
type DataPie struct {
	ID        int64
	Name      string
	ImageURL  string
	Slices    int
	Price     float64
	Purchases []*purchase.Purchase
	Labels    []*label.Label
}

// Pie includes DataPie, and any computed values
type Pie struct {
	DataPie
	RemainingSlices int
}
