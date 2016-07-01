package purchase

import "time"

const (
	// Key is the format by which purchase key names are made
	Key = "purchase_%d"

	// CountKey is the key use to store the id incrementer
	CountKey = "purchaseCounter"
)

// DataPurchase is the version of the purchase kept in the database
type DataPurchase struct {
	Customer  int64
	Amount    int
	Pie       int64
	Timestamp time.Time
}

// Purchase includes DataPurchase and any computed values
type Purchase struct {
	DataPurchase
}

// Populate fills in any computed values
func (p *Purchase) Populate() {
}
