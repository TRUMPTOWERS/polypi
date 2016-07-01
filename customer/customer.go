package customer

// DataCustomer is the customer that is stored in the database
type DataCustomer struct {
	ID       int64
	Username string
}

// Customer includes DataCustomer, and can also include computed values
type Customer struct {
	DataCustomer
}
