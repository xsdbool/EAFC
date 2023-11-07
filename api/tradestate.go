package api

// AuctionStatus represents the different auction statuses.
type AuctionStatus string

const (
	ACTIVE   AuctionStatus = "active"
	CLOSED   AuctionStatus = "closed"
	EXPIRED  AuctionStatus = "expired"
	INACTIVE AuctionStatus = "inactive"
	INVALID  AuctionStatus = "invalid"
	// Add more auction statuses as needed
)
