package main

// BidStatus represents the different bid statuses.
type BidStatus string

const (
	BUYNOW  BidStatus = "buyNow"
	HIGHEST BidStatus = "highest"
	NONE    BidStatus = "none"
	OUTBID  BidStatus = "outbid"
)
