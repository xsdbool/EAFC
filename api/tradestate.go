package api

// AuctionStatus represents the different auction statuses.
type AuctionStatus string

const (
	Active   AuctionStatus = "active"
	Closed   AuctionStatus = "closed"
	Expired  AuctionStatus = "expired"
	Inactive AuctionStatus = "inactive"
	Invalid  AuctionStatus = "invalid"
	// Add more auction statuses as needed
)
