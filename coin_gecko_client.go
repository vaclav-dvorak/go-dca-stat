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
)

type priceResponse struct {
	Prices [][]float64 `json:"prices"`
}

type stat struct {
	sum   float64
	count int64
	date  int
	avg   float64
	min   int
	min2  int
	score float64
}
type dcaStat struct {
	week  map[int]stat
	month map[int]stat
}

func getPriceData(cur string) (ret dcaStat, err error) {
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
	q.Add("days", strconv.Itoa(days+interval))
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

	// for i := 0; i < len(data.Prices)-1; i++ { //? last value is current price so we ignore it
	// 	tm := time.UnixMilli(int64(data.Prices[i][0]))
	// 	dow := int(tm.Weekday()) //* zero == sunday
	// 	if _, ok := ret.week[dow]; !ok {
	// 		ret.week[dow] = stat{sum: 0, count: 0, date: dow}
	// 	}
	// 	if entry, ok := ret.week[dow]; ok {
	// 		entry.sum += data.Prices[i][1]
	// 		entry.count++
	// 		ret.week[dow] = entry
	// 	}

	// 	day := int(tm.Day())
	// 	if _, ok := ret.month[day]; !ok {
	// 		ret.month[day] = stat{sum: 0, count: 0, date: day}
	// 	}
	// 	if entry, ok := ret.month[day]; ok {
	// 		entry.sum += data.Prices[i][1]
	// 		entry.count++
	// 		ret.month[day] = entry
	// 	}
	// }

	// // calc averages
	// for k, v := range ret.week {
	// 	v.avg = v.sum / float64(v.count)
	// 	ret.week[k] = v
	// }

	// for k, v := range ret.month {
	// 	v.avg = v.sum / float64(v.count)
	// 	ret.month[k] = v
	// }
	sumMin := make([]int, 7)
	sumMin2 := make([]int, 7)
	for n := 0; n < interval; n++ {
		min, min2, weekR := assessData(data.Prices[n : days+n])
		sumMin[min.date]++
		sumMin2[min2.date]++
		ret.week = weekR
	}

	for i, v := range ret.week {
		v.min = sumMin[ret.week[i].date]
		v.min2 = sumMin2[ret.week[i].date]
		ret.week[i] = v
	}

	return
}

func assessData(data [][]float64) (min, min2 stat, week map[int]stat) {
	week = make(map[int]stat, 0)
	for i := 0; i < len(data); i++ {
		tm := time.UnixMilli(int64(data[i][0]))
		dow := int(tm.Weekday()) //? zero == sunday
		if _, ok := week[dow]; !ok {
			week[dow] = stat{sum: 0, count: 0, date: dow}
		}
		if entry, ok := week[dow]; ok {
			entry.sum += data[i][1]
			entry.count++
			week[dow] = entry
		}
	}

	// calc averages
	mn, mn2 := 10000000.0, 10000000.0
	for k, v := range week {
		v.avg = v.sum / float64(v.count)
		week[k] = v
		if v.avg < mn {
			mn2 = mn
			mn = v.avg
			min2 = min
			min = week[k]
		} else if v.avg < mn2 {
			mn2 = v.avg
			min2 = week[k]
		}
	}
	return
}
