package api

// BidStatus represents the different bid statuses.
type BidStatus string

const (
	BuyNow  BidStatus = "buyNow"
	Highest BidStatus = "highest"
	None    BidStatus = "none"
	Outbid  BidStatus = "outbid"
)
