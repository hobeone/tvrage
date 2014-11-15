package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/drbig/tvrage"
)

const (
	VERSION = `0.0.1`
	TIMEFMT = `2006-01-02`
)

var (
	flagShows    bool
	flagEpisodes bool
	flagVersion  bool
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [option] show show...\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Without options prints last and next episode for the first matched show.\n")
		fmt.Fprintf(os.Stderr, "The -s, -e and -v options are mutually exclusive.\n")
		fmt.Fprintf(os.Stderr, "Actions described are executed for each show argument.\n\n")
		flag.PrintDefaults()
	}
	flag.BoolVar(&flagShows, "s", false, "print all matched shows")
	flag.BoolVar(&flagEpisodes, "e", false, "print all episodes for first matched show")
	flag.BoolVar(&flagVersion, "v", false, "print version")
}

func main() {
	flag.Parse()
	if flagVersion {
		fmt.Fprintf(os.Stderr, "tvrage command version: %s\n", VERSION)
		fmt.Fprintf(os.Stderr, "tvrage library version: %s\n", tvrage.VERSION)
		os.Exit(0)
	}
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(0)
	}

	for _, name := range flag.Args() {
		ss, err := tvrage.Search(name)
		if err != nil {
			fmt.Printf("%s - error: %s\n\n", name, err)
			continue
		}
		if flagShows {
			fmt.Printf("    %s\n", name)
			for idx, s := range ss {
				fmt.Printf("%2d. %s\n", idx+1, s)
			}
			fmt.Println()
			continue
		}

		fmt.Println(ss[0])
		es, err := tvrage.EpisodeList(ss[0].ID)
		if err != nil {
			fmt.Printf("error: %s\n\n", err)
			continue
		}
		if flagEpisodes {
			for idx, e := range es {
				fmt.Printf("%3d. %s\n", idx+1, e)
			}
			fmt.Println()
			continue
		}

		if ep, found := es.Last(); found {
			fmt.Printf("\tLAST: %s (%s, %s)\n", ep, ep.AirDate.Format(TIMEFMT), ep.DeltaDays())
		} else {
			fmt.Printf("\tLAST: Unknown\n")
		}
		if ep, found := es.Next(); found {
			fmt.Printf("\tNEXT: %s (%s, %s)\n", ep, ep.AirDate.Format(TIMEFMT), ep.DeltaDays())
		} else {
			fmt.Printf("\tNEXT: Unknown\n")
		}
		fmt.Println()
	}
}
