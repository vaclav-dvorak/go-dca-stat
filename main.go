package main

import (
	"fmt"
	"sort"

	"github.com/gosuri/uitable"
	log "github.com/sirupsen/logrus"
)

const (
	cur               = "czk"
	days              = 1820 //? 5*52*7 = 5 years
	interval          = 280
	scoreBase float64 = 100
	price     float64 = 100
	stprice   float64 = 98
	ndprice   float64 = 99
)

func calcScore(day stat) (score float64) {
	score = float64(day.min)*(scoreBase/stprice) + float64(day.min2)*(scoreBase/ndprice) + float64(interval-day.min2-day.min)*(scoreBase/price)
	return
}

func main() {
	dow := []string{0: "Nedele", 1: "Pondeli", 2: "Utery", 3: "Streda", 4: "Ctvrtek", 5: "Patek", 6: "Sobota"}
	printWelcome()
	dcaData, err := getPriceData(cur)
	if err != nil {
		log.Fatal("%v", err)
	}

	dcaWeekS := make([]stat, len(dcaData.week))
	for k := range dcaData.week {
		dcaWeekS[k] = dcaData.week[k]
		dcaWeekS[k].score = calcScore(dcaWeekS[k])
	}
	sort.Slice(dcaWeekS, func(i, j int) bool {
		return dcaWeekS[i].score > dcaWeekS[j].score
	})

	table := uitable.New()
	table.AddRow("day", "price", "(%)", "min stat", "score")
	table.AddRow("===", "=====", "===", "========", "=====")
	for _, v := range dcaWeekS {
		p := (dcaWeekS[3].avg - v.avg) / dcaWeekS[3].avg
		table.AddRow(dow[v.date], fmt.Sprintf("%0.2f %s", v.avg, cur), fmt.Sprintf("%0.3f", p*100), fmt.Sprintf("%d/%d/(%d+%d)", interval, v.min+v.min2, v.min, v.min2), fmt.Sprintf("%0.3f", v.score))
	}
	fmt.Println(table)

}
