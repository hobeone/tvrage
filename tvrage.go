package tvrage

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Show struct {
	ID             int      `xml:"showid"`
	Name           string   `xml:"name"`
	Link           string   `xml:"link"`
	Country        string   `xml:"country"`
	Started        int      `xml:"started"`
	Ended          int      `xml:"ended"`
	Seasons        int      `xml:"seasons"`
	Status         string   `xml:"status"`
	Classification string   `xml:"classification"`
	Genres         []string `xml:"genres>genre"`
}

func (s Show) String() string {
	return fmt.Sprintf("%s [%d - %s]", s.Name, s.Started, s.Status)
}

type tvrageTime struct {
	time.Time
}

func (t *tvrageTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parsed, err := time.Parse(TIMEFMT, v)
	if err != nil {
		return nil
	}
	*t = tvrageTime{parsed}
	return nil
}

type Episode struct {
	Season     int
	Ordinal    int        `xml:"epnum"`
	Number     int        `xml:"seasonnum"`
	Production string     `xml:"prodnum"`
	AirDate    tvrageTime `xml:"airdate"`
	Link       string     `xml:"link"`
	Title      string     `xml:"title"`
}

func (e Episode) String() string {
	return fmt.Sprintf(`S%02dE%02d "%s"`, e.Season, e.Number, e.Title)
}

func (e *Episode) DeltaDays() string {
	d := int(e.AirDate.Sub(time.Now()).Hours() / 24.0)
	if d < 0 {
		if d == -1 {
			return "yesterday"
		} else {
			return fmt.Sprintf("%d days ago", -d)
		}
	} else if d > 0 {
		if d == 1 {
			return "tomorrow"
		} else {
			return fmt.Sprintf("in %d days", d+1)
		}
	} else {
		return "today"
	}
}

type Episodes []Episode

func (es Episodes) Last() (Episode, bool) {
	var r Episode
	t := time.Now()
	for _, e := range es {
		if e.AirDate.IsZero() {
			continue
		}
		if e.AirDate.Before(t) {
			r = e
		}
	}
	if r.AirDate.IsZero() {
		return r, false
	} else {
		return r, true
	}
}

func (es Episodes) Next() (Episode, bool) {
	var r Episode
	t := time.Now()
	for _, e := range es {
		if e.AirDate.IsZero() {
			continue
		}
		if e.AirDate.After(t) {
			return e, true
		}
	}
	return r, false
}

type resultSeason struct {
	Number   int       `xml:"no,attr"`
	Episodes []Episode `xml:"episode"`
}

type resultEpisodeList struct {
	Total   int            `xml:"totalseasons"`
	Seasons []resultSeason `xml:"Episodelist>Season"`
}

type resultSearch struct {
	Shows []Show `xml:"show"`
}

const (
	SEARCHURL = `http://services.tvrage.com/feeds/search.php?show=%s`      // URL for show searching
	EPLISTURL = `http://services.tvrage.com/feeds/episode_list.php?sid=%d` // URL for episode list
	TIMEFMT   = `2006-01-02`                                               // time.Parse format string for air date
)

var (
	Client = &http.Client{} // default HTTP client
)

func parseSearchResult(in io.Reader) ([]Show, error) {
	r := resultSearch{}
	x := xml.NewDecoder(in)
	if err := x.Decode(&r); err != nil {
		return nil, err
	}
	return r.Shows, nil
}

func Search(name string) ([]Show, error) {
	q := fmt.Sprintf(SEARCHURL, url.QueryEscape(name))
	r, err := Client.Get(q)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return parseSearchResult(r.Body)
}

func parseEpisodeListResult(in io.Reader) (Episodes, error) {
	var es Episodes
	r := resultEpisodeList{}
	x := xml.NewDecoder(in)
	if err := x.Decode(&r); err != nil {
		return nil, err
	}
	for _, s := range r.Seasons {
		for _, e := range s.Episodes {
			e.Season = s.Number
			es = append(es, e)
		}
	}
	return es, nil
}

func EpisodeList(id int) (Episodes, error) {
	q := fmt.Sprintf(EPLISTURL, id)
	r, err := Client.Get(q)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return parseEpisodeListResult(r.Body)
}
