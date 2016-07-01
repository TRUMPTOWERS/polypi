package customer

const (
	// Key is the format by which the keys for customers are formed
	Key = "customer_%d"
	// IDByName is a hash index for getting IDs via names
	IDByName = "customerName_%s"

	// CountKey is the key use to store the id incrementer
	CountKey = "customerCounter"
)

// DataCustomer is the customer that is stored in the database
type DataCustomer struct {
	ID       int64
	Username string
}

// Customer includes DataCustomer, and can also include computed values
type Customer struct {
	DataCustomer
}

// Populate fills in any computed values
func (c *Customer) Populate() {
}
