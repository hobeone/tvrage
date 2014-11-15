// See LICENSE.txt for licensing information.

package tvrage

import (
	"strings"
	"testing"
	"time"
)

func TestParseSearchResult(t *testing.T) {
	input := `<Results>
<show>
<showid>5410</showid>
<name>Supernatural</name>
<link>http://www.tvrage.com/Supernatural</link>
<country>US</country>
<started>2005</started>
<ended>0</ended>
<seasons>10</seasons>
<status>Returning Series</status>
<classification>Scripted</classification>
<genres>
<genre>Action</genre>
<genre>Adventure</genre>
<genre>Drama</genre>
<genre>Horror/Supernatural</genre>
<genre>Sci-Fi</genre>
</genres>
</show>
<show>
<showid>2032</showid>
<name>Supernatural (1977)</name>
<link>http://www.tvrage.com/supernatural-1977</link>
<country>UK</country>
<started>1977</started>
<ended>1977</ended>
<seasons>1</seasons>
<status>Canceled/Ended</status>
<classification>0</classification>
<genres/>
</show>
<show>
<showid>27870</showid>
<name>Supernatural: The Animation</name>
<link>http://www.tvrage.com/supernatural-the-animation</link>
<country>JP</country>
<started>2011</started>
<ended>2011</ended>
<seasons>1</seasons>
<status>Canceled/Ended</status>
<classification>Animation</classification>
<genres>
<genre>Anime</genre>
<genre>Horror/Supernatural</genre>
</genres>
</show>
</Results>`

	res, err := parseSearchResult(strings.NewReader(input))
	if err != nil {
		t.Errorf("Decode error: %s", err)
		t.FailNow()
	}

	if len(res) != 3 {
		t.Errorf("Length mismatch: %d != %d", len(res), 3)
	}
	if res[0].ID != 5410 {
		t.Errorf("(1) Show ID mismatch: %d != %d", res[0].ID, 5410)
	}
	if res[1].Name != "Supernatural (1977)" {
		t.Errorf("(2) Show Name mismatch: %s != %s", res[1].Name, "Supernatural (1977)")
	}
	if len(res[2].Genres) != 2 {
		t.Errorf("(3) Genres length mismatch: %d != %d", len(res[2].Genres), 2)
	}

	if res[0].String() != "Supernatural [2005 - Returning Series]" {
		t.Errorf("Show stringer output mismatch")
	}

	res, err = parseSearchResult(strings.NewReader(``))
	if err == nil {
		t.Errorf("Didn't fail with empty data")
		t.FailNow()
	}
}

func TestSearchLive(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	r, err := Search("archer (2009)")
	if err != nil {
		t.Errorf("Search error: %s", err)
		t.FailNow()
	}

	if len(r) < 1 {
		t.Errorf("Less than one Show decoded")
	}
}

func TestParseEpisodeListResult(t *testing.T) {
	input := `<Show>
<name>Supernatural</name>
<totalseasons>10</totalseasons>
<Episodelist>
<Season no="1">
<episode>
<epnum>0</epnum>
<seasonnum>00</seasonnum>
<prodnum/>
<airdate>0000-00-00</airdate>
<link>
http://www.tvrage.com/Supernatural/episodes/1065190732
</link>
<title>Unaired Pilot</title>
</episode>
<episode>
<epnum>1</epnum>
<seasonnum>01</seasonnum>
<prodnum>475285</prodnum>
<airdate>2005-09-13</airdate>
<link>http://www.tvrage.com/Supernatural/episodes/166205</link>
<title>Pilot</title>
</episode>
</Season>
<Season no="2">
<episode>
<epnum>23</epnum>
<seasonnum>01</seasonnum>
<prodnum>3T5501</prodnum>
<airdate>2006-09-28</airdate>
<link>http://www.tvrage.com/Supernatural/episodes/386441</link>
<title>In My Time of Dying</title>
</episode>
<episode>
<epnum>24</epnum>
<seasonnum>02</seasonnum>
<prodnum>3T5502</prodnum>
<airdate>2006-10-05</airdate>
<link>http://www.tvrage.com/Supernatural/episodes/412873</link>
<title>Everybody Loves a Clown</title>
</episode>
<episode>
<epnum>25</epnum>
<seasonnum>43</seasonnum>
<prodnum>3T5502</prodnum>
<airdate>2022-12-15</airdate>
<link>http://www.tvrage.com/Supernatural/episodes/412873</link>
<title>Test Fake Episode</title>
</episode>
<episode>
<epnum>26</epnum>
<seasonnum>47</seasonnum>
<prodnum>3T5502</prodnum>
<airdate>2022-12-32</airdate>
<link>http://www.tvrage.com/Supernatural/episodes/412873</link>
<title>Test Fake Episode Two</title>
</episode>

</Season>
</Episodelist>
</Show>`

	res, err := parseEpisodeListResult(strings.NewReader(input))
	if err != nil {
		t.Errorf("Decode error: %s", err)
		t.FailNow()
	}

	if len(res) != 6 {
		t.Errorf("Length mismatch: %d != %d", len(res), 4)
	}
	if res[0].Title != "Unaired Pilot" {
		t.Errorf("(1) Episode Title mismatch: %s != %s", res[0].Title, "Unaired Pilot")
	}
	if res[1].Season != 1 && res[1].Number != 1 {
		t.Errorf("(2) Episode Season/Number mismatch: S%dE%d != %s", res[1].Season, res[1].Number, "S1E1")
	}
	testDate, err := time.Parse(TIMEFMT, "2006-09-28")
	if err != nil {
		t.Errorf("Couldn't make testDate: %s", err)
	} else {
		if !res[2].AirDate.Equal(testDate) {
			t.Errorf("(3) AirDate mismatch: %s != %s", res[2].AirDate, testDate)
		}
	}

	if res[2].String() != `S02E01 "In My Time of Dying"` {
		t.Errorf("Episode stringer output mismatch")
	}

	lep, found := res.Last()
	if !found {
		t.Errorf("Didn't find the last episode")
	} else {
		if lep.String() != `S02E02 "Everybody Loves a Clown"` {
			t.Errorf("Wrong episode returned by Last(): %s", lep)
		}
	}

	nep, found := res.Next()
	if !found {
		t.Errorf("Didn't find the next episode")
	} else {
		if nep.String() != `S02E43 "Test Fake Episode"` {
			t.Errorf("Wrong episode returned by Next(): %s", nep)
		}
	}

	res, err = parseEpisodeListResult(strings.NewReader(``))
	if err == nil {
		t.Errorf("Didn't fail with empty data")
		t.FailNow()
	}
}

func TestEpisodeListLive(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	r, err := EpisodeList(5410)
	if err != nil {
		t.Errorf("EpisodeList error: %s", err)
		t.FailNow()
	}

	if len(r) < 1 {
		t.Errorf("Less than one Episode decoded")
	}
}

type deltaTest struct {
	episode Episode
	output  string
}

func TestDeltaDays(t *testing.T) {
	r := time.Now()
	cases := []deltaTest{
		deltaTest{Episode{AirDate: tvrageTime{r.Add(12 * time.Hour)}}, "today"},
		deltaTest{Episode{AirDate: tvrageTime{r.Add(25 * time.Hour)}}, "tomorrow"},
		deltaTest{Episode{AirDate: tvrageTime{r.Add(24 * 3 * time.Hour)}}, "in 3 days"},
		deltaTest{Episode{AirDate: tvrageTime{r.Add(24 * 43 * time.Hour)}}, "in 43 days"},
		deltaTest{Episode{AirDate: tvrageTime{r.Add(-25 * time.Hour)}}, "yesterday"},
		deltaTest{Episode{AirDate: tvrageTime{r.Add(-24 * 3 * time.Hour)}}, "3 days ago"},
		deltaTest{Episode{AirDate: tvrageTime{r.Add(-24 * 120 * time.Hour)}}, "120 days ago"},
	}

	for idx, c := range cases {
		if c.episode.DeltaDays() != c.output {
			t.Errorf("(%d) Mismatch: %s != %s", idx+1, c.episode.DeltaDays(), c.output)
		}
	}
}
