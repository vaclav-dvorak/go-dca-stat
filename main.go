package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"
	log "github.com/sirupsen/logrus"
)

const (
	cur     = "czk"
	daysNum = 1820 //? 5*52*7 = 5 years
	// daysNum = 1092 //? 3*52*7 = 3 years
	// daysNum           = 364 //? 1*52*7 = 1 years
)

func main() {
	dow := []string{0: "Nedele", 1: "Pondeli", 2: "Utery", 3: "Streda", 4: "Ctvrtek", 5: "Patek", 6: "Sobota"}
	printWelcome()
	dcaData, err := getPriceData(cur, daysNum)
	if err != nil {
		log.Fatalf("+%v", err)
	}

	dcaWeekS := make([]stat, len(dcaData.week))
	for k := range dcaData.week {
		dcaWeekS[k] = dcaData.week[k]
	}
	sort.Slice(dcaWeekS, func(i, j int) bool {
		return dcaWeekS[i].score > dcaWeekS[j].score
	})

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"day", "price", "count", "score"})
	for _, v := range dcaWeekS {
		t.AppendRow(table.Row{
			dow[v.date], fmt.Sprintf("%0.0f %s", v.avg, cur), v.count, fmt.Sprintf("%0.3f", v.score*100),
		})
	}
	t.SetStyle(table.StyleLight)
	t.Render()
}
