package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

type EaFcAuthedclient struct {
	maxTimeout int32
	minTimeout int32
	sessionID  string
	client     http.Client
}

func NewEAFCAuthedClient(sessionID string, maxTimeout int32, minTimeout int32) *EaFcAuthedclient {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		panic("cookie jar fail")
	}

	return &EaFcAuthedclient{
		maxTimeout: maxTimeout,
		minTimeout: minTimeout,
		sessionID:  sessionID,
		client: http.Client{
			Jar: jar,
		},
	}
}

func (c *EaFcAuthedclient) Sleep() {
	random_timeout := rand.Int31n(c.maxTimeout - c.minTimeout)
	timeout := random_timeout + int32(c.minTimeout)
	time.Sleep(time.Duration(time.Duration(timeout) * time.Millisecond))
}

func (c *EaFcAuthedclient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-UT-SID", c.sessionID)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[%d] %s %s\n", resp.StatusCode, resp.Request.URL.Path, resp.Request.URL.RawQuery)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[%d] %s\n", resp.StatusCode, UtasErrorCode[resp.StatusCode])

		var v any
		err := json.NewDecoder(resp.Body).Decode(&v)
		if err == nil {
			fmt.Printf("%v\n", v)
		}
		if err != nil {
			body, _ := io.ReadAll(resp.Body)
			fmt.Println(string(body))
		}
	}
	switch code := resp.StatusCode; code {
	case http.StatusUnauthorized:
		panic("session expired")
		//return nil, ErrSessionExpired
	case 458:
		return nil, ErrCaptchaRequired
	}
	return resp, err
}

func (c *EaFcAuthedclient) extractJSON(req *http.Request, v any) (*http.Response, error) {
	response, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *EaFcAuthedclient) SearchTransfermarket(query *url.Values) (*TransfermarketResponse, error) {
	url, err := url.Parse("https://utas.mob.v2.fut.ea.com/ut/game/fc24/transfermarket")
	if err != nil {
		return nil, err
	}
	url.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	response, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var resp TransfermarketResponse
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil

}

func (c *EaFcAuthedclient) Bid(bidPrice int, auction AuctionInfo) (*BidResponse, error) {

	if bidPrice <= auction.CurrentBid {
		fmt.Println("bidPrice <= CurrentBid")
		return nil, ErrBidding
	}

	if auction.TradeState == string(Closed) {
		fmt.Println("trade closed")
		return nil, ErrBidding
	}

	// Construct the PUT request URL with the TradeId
	bidURL := fmt.Sprintf("https://utas.mob.v2.fut.ea.com/ut/game/fc24/trade/%d/bid", auction.TradeId)

	// Create a JSON request body with the bid price
	payload := map[string]int{"bid": bidPrice}
	jsonpayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, bidURL, bytes.NewBuffer(jsonpayload))
	if err != nil {
		return nil, err
	}

	var bidResponse BidResponse
	response, err := c.extractJSON(req, &bidResponse)
	if err != nil {
		return nil, err
	}

	switch code := response.StatusCode; code {
	case http.StatusOK:
		return &bidResponse, nil
	case 512:
		fallthrough
	case 521:
		fallthrough
	case 426:
		return nil, ErrFunctionDisabled
	case 461:
		body, _ := io.ReadAll(response.Body)
		if strings.Contains(string(body), "Watchlist is full") {
			return nil, ErrWatchlistIsFull
		}
	}
	return nil, ErrBidding
}

func (c *EaFcAuthedclient) Buy(auction AuctionInfo) (*BidResponse, error) {

	if auction.TradeState == string(Closed) {
		return nil, ErrBidding
	}

	// Construct the PUT request URL with the TradeId
	bidURL := fmt.Sprintf("https://utas.mob.v2.fut.ea.com/ut/game/fc24/trade/%d/bid", auction.TradeId)

	// Create a JSON request body with the bid price
	payload := map[string]int{
		"bid": auction.BuyNowPrice,
	}
	jsonpayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, bidURL, bytes.NewBuffer(jsonpayload))
	if err != nil {
		return nil, err
	}
	var bidResponse BidResponse
	response, err := c.extractJSON(req, &bidResponse)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, ErrBidding
	}
	return &bidResponse, nil
}

func (c *EaFcAuthedclient) ListAuction(auction AuctionInfo, buyNowPrice int, startingBid int, duration int) (*AuctionHouseResponse, error) {
	url := "https://utas.mob.v2.fut.ea.com/ut/game/fc24/auctionhouse"

	// Define the payload data as a Go struct
	payload := struct {
		BuyNowPrice int `json:"buyNowPrice"`
		Duration    int `json:"duration"`
		ItemData    struct {
			ID int `json:"id"`
		} `json:"itemData"`
		StartingBid int `json:"startingBid"`
	}{
		BuyNowPrice: buyNowPrice,
		Duration:    duration,
		ItemData: struct {
			ID int `json:"id"`
		}{
			ID: auction.ItemData.ID,
		},
		StartingBid: startingBid,
	}

	jsonpayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonpayload))
	if err != nil {
		return nil, err
	}

	var auctionHouseResponse AuctionHouseResponse
	response, err := c.extractJSON(req, &auctionHouseResponse)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, ErrListing
	}
	return &auctionHouseResponse, nil
}

func (c *EaFcAuthedclient) Relist() (*RelistResponse, error) {
	url := "https://utas.mob.v2.fut.ea.com/ut/game/fc24/auctionhouse/relist"

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return nil, err
	}

	var reslistResponse RelistResponse
	response, err := c.extractJSON(req, &reslistResponse)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, ErrRelist
	}
	return &reslistResponse, nil

}

func (c *EaFcAuthedclient) ItemToPile(pile string, auction AuctionInfo) (*ItemResponse, error) {
	url := "https://utas.mob.v2.fut.ea.com/ut/game/fc24/item"

	// Define the payload data
	payload := map[string]interface{}{
		"itemData": []map[string]interface{}{
			{
				"id":      auction.ItemData.ID,
				"pile":    pile,
				"tradeId": auction.TradeId,
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	var itemResponse ItemResponse
	response, err := c.extractJSON(req, &itemResponse)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, ErrMoveItem
	}
	return &itemResponse, nil
}

func (c *EaFcAuthedclient) Watchlist() (*WatchlistResponse, error) {
	// Construct the DELETE request URL with the list of tradeIds
	watchlistURL := "https://utas.mob.v2.fut.ea.com/ut/game/fc24/watchlist"

	req, err := http.NewRequest(http.MethodGet, watchlistURL, nil)
	if err != nil {
		return nil, err
	}

	var watchlistResponse WatchlistResponse
	response, err := c.extractJSON(req, &watchlistResponse)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, ErrWatchlist
	}
	return &watchlistResponse, nil

}

func (c *EaFcAuthedclient) ClearWatchlist() error {
	wr, err := c.Watchlist()
	if err != nil {
		return ErrClearWatchlist
	}

	c.Sleep()
	var tradeIds []int
	for _, auction := range wr.AuctionInfo {
		//if auction.BidState == string(Outbid) && auction.TradeState == string(Closed) {
		if auction.BidState == string(Outbid) {
			tradeIds = append(tradeIds, auction.TradeId)
		}
	}
	if len(tradeIds) == 0 {
		return ErrClearWatchlist
	}
	// Construct the DELETE request URL with the list of tradeIds
	watchlistURL := "https://utas.mob.v2.fut.ea.com/ut/game/fc24/watchlist?tradeId="

	// Append the tradeIds to the URL
	for i, tradeID := range tradeIds {
		if i > 0 {
			watchlistURL += ","
		}
		watchlistURL += fmt.Sprintf("%d", tradeID)
	}

	req, err := http.NewRequest(http.MethodDelete, watchlistURL, nil)
	if err != nil {
		return err
	}

	response, err := c.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return ErrClearWatchlist
	}
	return nil
}

func (c *EaFcAuthedclient) ClearSold() error {
	// Construct the DELETE request URL with the list of tradeIds
	url := "https://utas.mob.v2.fut.ea.com/ut/game/fc24/trade/sold"

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	response, err := c.Do(req)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return ErrClearSold
	}
	return nil
}
