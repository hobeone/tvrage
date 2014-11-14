package main

import (
	"fmt"
	"os"

	"github.com/drbig/tvrage"
)

func main() {
	for _, name := range os.Args[1:] {
		s, err := tvrage.Search(name)
		if err != nil {
			fmt.Printf("%s - error: %s\n\n", name, err)
			continue
		}
		fmt.Println(s[0])
		es, err := tvrage.EpisodeList(s[0].ID)
		if err != nil {
			fmt.Printf("error: %s\n\n", err)
			continue
		}

		if lep, found := es.Last(); found {
			fmt.Printf("\tLAST: %s (%s, %s)\n", lep, lep.AirDate.Format(tvrage.TIMEFMT), lep.DeltaDays())
		} else {
			fmt.Printf("\tLAST: Unknown\n")
		}

		if nep, found := es.Next(); found {
			fmt.Printf("\tNEXT: %s (%s, %s)\n", nep, nep.AirDate.Format(tvrage.TIMEFMT), nep.DeltaDays())
		} else {
			fmt.Printf("\tNEXT: Unknown\n")
		}

		fmt.Println()
	}
}
