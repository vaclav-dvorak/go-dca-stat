package main

import (
	"fmt"
)

var (
	logoSmall = []string{
		`       ▄ ▄   `,
		`      ▀█▀▀▀▀▄`,
		`       █▄▄▄▄▀`,
		`       █    █`,
		`      ▀▀█▀█▀ `,
	}
	version, date = "(devel)", "now"
)

func printWelcome() {
	for i := 0; i < len(logoSmall); i++ {
		fmt.Printf("%s%s%s", orange, logoSmall[i], reset)
		if i < (len(logoSmall) - 1) {
			fmt.Print("\n")
		}
	}
	fmt.Printf(" %s%s @%s%s // dollar cost averaging stats\n", green, version, date, reset)
	fmt.Printf("all price data are powered by CoinGecko: %shttps://www.coingecko.com/%s\n\n", blue, reset)
}
