package main

import "errors"

var (
	ErrSessionExpired   = errors.New("session expired or unauthorized")
	ErrFunctionDisabled = errors.New("function currently disabled")
	ErrBidding          = errors.New("error bidding on item")
	ErrListing          = errors.New("error listing item")
	ErrWatchlist        = errors.New("error getting watchlist")
	ErrClearWatchlist   = errors.New("error clearing watchlist")
	ErrReslist          = errors.New("error reslisting items")
	ErrClearSold        = errors.New("cannot clear sold items")
	ErrHTTPRequest      = errors.New("HTTP request error")
	ErrInvalidBidPrice  = errors.New("bid price is invalid")
	ErrMoveItem         = errors.New("error putting Item to pile")
)
