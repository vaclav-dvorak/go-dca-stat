// Package main implements whole functionality of this tool
package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	timeoutSec = 5
	buyAmount  = 1000 //? We spent amount of currency on every buy
)

type priceResponse struct {
	Prices [][]float64 `json:"prices"`
}

type stat struct {
	sum   float64
	count int64
	avg   float64
	date  int
	score float64
}
type dcaStat struct {
	week  map[int]stat
	month map[int]stat
}

func getPriceData(cur string, days int) (ret dcaStat, err error) {
	var (
		data priceResponse
	)
	ret.week = make(map[int]stat, 0)
	ret.month = make(map[int]stat, 0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*timeoutSec))
	defer cancel()
	header := http.Header{
		"Accept":       []string{"application/json"},
		"Content-Type": []string{"application/json"},
	}
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.coingecko.com/api/v3/coins/bitcoin/market_chart", nil)
	req.Header = header
	q := req.URL.Query()
	q.Add("days", strconv.Itoa(days))
	q.Add("vs_currency", cur)
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return
	}

	body, _ := io.ReadAll(res.Body)
	defer func() {
		_ = res.Body.Close()
	}()
	if err = json.Unmarshal(body, &data); err != nil {
		return
	}

	week := make(map[int]stat, 0)
	for i := 0; i < len(data.Prices)-1; i++ {
		d := data.Prices[i]
		tm := time.UnixMilli(int64(d[0]))
		dow := int(tm.Weekday()) //? zero == sunday
		if _, ok := week[dow]; !ok {
			week[dow] = stat{date: dow}
		}
		if entry, ok := week[dow]; ok {
			entry.sum += d[1]
			entry.count++
			entry.score += buyAmount / d[1]
			entry.avg = entry.sum / float64(entry.count)
			week[dow] = entry
		}
	}

	ret.week = week
	return
}
