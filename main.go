package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"
	log "github.com/sirupsen/logrus"
)

const (
	cur        = "czk"
	yearNum    = 5
	beautifier = 1000
)

func main() {
	dow := []string{0: "Nedele", 1: "Pondeli", 2: "Utery", 3: "Streda", 4: "Ctvrtek", 5: "Patek", 6: "Sobota"}
	printWelcome()
	years := make([]dcaStat, yearNum)
	for year := 0; year < yearNum; year++ {
		log.Infof("Fetching year (%d/%d)...\n", year+1, yearNum)
		dcaData, err := getPriceData(cur, (year+1)*364) //? 364 = 1*52*7 = 1 year
		if err != nil {
			log.Fatalf("+%v", err)
		}
		years[year] = dcaData
	}
	var dcaWeekS []stat
	dcaWeekS = append(dcaWeekS, years[0].week...) // we need deep copy of slice
	sort.Slice(dcaWeekS, func(i, j int) bool {
		return dcaWeekS[i].score > dcaWeekS[j].score
	})

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"day", "price", "count", "score (5y)", "score (3y)", "score (1y)"})
	for _, v := range dcaWeekS {
		t.AppendRow(table.Row{
			dow[v.date], fmt.Sprintf("%0.0f %s", v.avg, cur), v.count, fmt.Sprintf("%0.3f", years[4].week[v.date].score*beautifier), fmt.Sprintf("%0.3f", years[2].week[v.date].score*beautifier), fmt.Sprintf("%0.3f", years[0].week[v.date].score*beautifier),
		})
	}
	t.SetStyle(table.StyleLight)
	t.Render()
}
