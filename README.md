# tvrage [![Build Status](https://travis-ci.org/drbig/tvrage.svg?branch=master)](https://travis-ci.org/drbig/tvrage) [![Coverage Status](https://img.shields.io/coveralls/drbig/tvrage.svg)](https://coveralls.io/r/drbig/tvrage?branch=master) [![GoDoc](https://godoc.org/github.com/drbig/tvrage?status.svg)](http://godoc.org/github.com/drbig/tvrage)

Package tvrage provides basic access to tvrage.com services for finding out the last and next episodes of a given TV show (plus a bit more), no API key required.

## Showcase

Using the included *demo* `tvrage`:

    $ ./tvrage
    Usage: ./tvrage [option] show show...
    
    Without options prints last and next episode for the first matched show.
    The -s, -e and -v options are mutually exclusive.
    Actions described are executed for each show argument.
    
      -e=false: print all episodes for first matched show
      -s=false: print all matched shows
      -v=false: print version

- - -

    $ ./tvrage -s archer
        archer
     1. Archer (2009) [2009 - Returning Series]
     2. Archer (1975) [1975 - Canceled/Ended]
     3. Meet Corliss Archer [1951 - Canceled/Ended]
     4. Elite Archery's Respect the Game TV [2011 - Returning Series]
     5. Cabela's American Archer [2013 - Returning Series]
     6. Mathew's Dominant Bucks [2010 - Returning Series]
     7. Archer's Goon [1992 - Canceled/Ended]
     8. The Sex Researchers [2011 - New Series]
     9. Archer's Choice with Ralph & Vicki [2001 - Returning Series]
    

- - -

    $ ./tvrage -e "game of thrones"
    Game of Thrones [2010 - Returning Series]
      1. S01E01 "Winter is Coming"
      2. S01E02 "The Kingsroad"
      3. S01E03 "Lord Snow"
      4. S01E04 "Cripples, Bastards, and Broken Things"
      5. S01E05 "The Wolf and the Lion"
      6. S01E06 "A Golden Crown"
      7. S01E07 "You Win or You Die"
      8. S01E08 "The Pointy End"
      9. S01E09 "Baelor"
     10. S01E10 "Fire and Blood"
     11. S02E01 "The North Remembers"
     12. S02E02 "The Night Lands"
     13. S02E03 "What is Dead May Never Die"
     14. S02E04 "Garden of Bones"
     15. S02E05 "The Ghost of Harrenhal"
     16. S02E06 "The Old Gods and the New"
     17. S02E07 "A Man Without Honor"
     18. S02E08 "The Prince of Winterfell"
     19. S02E09 "Blackwater"
     20. S02E10 "Valar Morghulis"
     21. S03E01 "Valar Dohaeris"
     22. S03E02 "Dark Wings, Dark Words"
     23. S03E03 "Walk of Punishment"
     24. S03E04 "And Now His Watch is Ended"
     25. S03E05 "Kissed by Fire"
     26. S03E06 "The Climb"
     27. S03E07 "The Bear and the Maiden Fair"
     28. S03E08 "Second Sons"
     29. S03E09 "The Rains of Castamere"
     30. S03E10 "Mhysa"
     31. S04E01 "Two Swords"
     32. S04E02 "The Lion and the Rose"
     33. S04E03 "Breaker of Chains"
     34. S04E04 "Oathkeeper"
     35. S04E05 "First of His Name"
     36. S04E06 "The Laws of Gods and Men"
     37. S04E07 "Mockingbird"
     38. S04E08 "The Mountain and the Viper"
     39. S04E09 "The Watchers on the Wall"
     40. S04E10 "The Children"
     41. S05E01 "Season 5, Episode 1"
    

- - -

    $ ./tvrage "sons of anarchy" "true dedective" vikings
    Sons of Anarchy [2008 - Final Season]
            LAST: S07E10 "Faith and Despondency" (2014-11-11, 4 days ago)
            NEXT: S07E11 "Suits of Woe" (2014-11-18, in 3 days)
    
    True Detective [2014 - Returning Series]
            LAST: S01E08 "Form and Void" (2014-03-09, 251 days ago)
            NEXT: Unknown
    
    Vikings [2013 - Returning Series]
            LAST: S02E10 "The Lord's Prayer" (2014-05-01, 198 days ago)
            NEXT: Unknown
    

## Contributing

Follow the usual GitHub development model:

1. Clone the repository
2. Make your changes on a separate branch
3. Make sure you run `gofmt` and `go test` before committing
4. Make a pull request

See licensing for legalese.

## Licensing

Standard two-clause BSD license, see LICENSE.txt for details.

Any contributions will be licensed under the same conditions.

Copyright (c) 2014 Piotr S. Staszewski
