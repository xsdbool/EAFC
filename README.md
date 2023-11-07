# EAFC
EA Sports Football Club


basic example:

```go
package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/xsdbool/EAFC/api"
)

func main() {
	const sid = "paste-ur-session-here"

	c := api.NewEAFCAuthedClient(sid, 1001, 1000)
  c.Relist()
  c.ClearSold()

  q := url.Values{}
	q.Set("lev", "gold")
	q.Set("type", "clubInfo")
	q.Set("cat", "kit")
	q.Set("authenticity", "true")
	q.Set("start", fmt.Sprintf("%d", 0))
	q.Set("num", fmt.Sprintf("%d", 21))
	q.Set("rarityIds", fmt.Sprintf("%d", 1))
	q.Set("maxb", fmt.Sprintf("%d", 450))

	auctions, err := client.SearchTransfermarket(&q)
```
