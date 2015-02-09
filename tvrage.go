// See LICENSE.txt for licensing information.

// Package tvrage provides basic access to tvrage.com services for finding out the last
// and next episodes of a given TV show (plus a bit more), no API key required.
package tvrage

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Show maps all available show data, as retrieved via Detailed Search and ShowInfo.
// Some fields come back with different names depending on the api used.
// Link, Country, Started and Ended are always filled out.
// ShowLink, OriginCountry, and StartDate may not be.
type Show struct {
	ID             int64      `xml:"showid"`
	Name           string     `xml:"name"`
	ShowName       string     `xml:"showname,omitempty"`
	Link           string     `xml:"link"`
	ShowLink       string     `xml:"showlink,omitempty"`
	Country        string     `xml:"country"`
	OriginCountry  string     `xml:"origin_country,omitempty"`
	Started        tvrageYear `xml:"started"`
	StartDate      tvrageYear `xml:"startdate"`
	Ended          tvrageYear `xml:"ended"`
	Seasons        int64      `xml:"seasons"`
	Status         string     `xml:"status"`
	Runtime        int64      `xml:"runtime"`
	Classification string     `xml:"classification"`
	Genres         []string   `xml:"genres>genre"`
	Network        string     `xml:"network"`
	Airtime        string     `xml:"airtime"`
	Airday         string     `xml:"airday"`
	Timezone       string     `xml:"timezone"`
}

// String returns a pretty string for a given Show.
func (s Show) String() string {
	return fmt.Sprintf("%s [%d - %s]", s.Name, s.Started, s.Status)
}

// tvrageTime is a thin shim over time.Time used to implement XML unmarshaling.
type tvrageTime struct {
	time.Time
}

// UnmarshalXML implements time.Time XML unmarshaling for tvrage.com air date format.
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

// tvrageYear is a shim over crappy date formats in the TVRage results.
// It can handle bare Years or the Mar/10/1997 format.  It will only have the year set in the result though.
type tvrageYear int64

// UnmarshalXML implements time.Time XML unmarshaling for year only tvrage.com formats.
func (t *tvrageYear) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	if strings.Contains(v, "/") {
		parts := strings.Split(v, "/")
		v = parts[len(parts)-1]
	}
	year, err := strconv.Atoi(v)
	if err != nil {
		return nil
	}
	*t = tvrageYear(year)
	return nil
}

// Episode maps all available episode data, as retrieved via EpisodeList.
type Episode struct {
	Season     int
	Ordinal    int        `xml:"epnum"`
	Number     int        `xml:"seasonnum"`
	Production string     `xml:"prodnum"`
	AirDate    tvrageTime `xml:"airdate"`
	Link       string     `xml:"link"`
	Title      string     `xml:"title"`
}

// String returns a pretty string for a given Episode.
func (e Episode) String() string {
	return fmt.Sprintf(`S%02dE%02d "%s"`, e.Season, e.Number, e.Title)
}

// DeltaDaysInt returns number of days since or to the episode's air date.
func (e *Episode) DeltaDaysInt() int {
	d := int(e.AirDate.Sub(time.Now()).Hours() / 24.0)
	if d > 1 {
		d++
	}
	return d
}

// DeltaDays returns a pretty string indicating the delta in days between now
// and the episode air date.
func (e *Episode) DeltaDays() string {
	d := e.DeltaDaysInt()
	if d < 0 {
		if d == -1 {
			return "yesterday"
		}
		return fmt.Sprintf("%d days ago", -d)
	} else if d > 0 {
		if d == 1 {
			return "tomorrow"
		}
		return fmt.Sprintf("in %d days", d)
	} else {
		return "today"
	}
}

// Episodes is a thin shim over []Episodes to enable methods on Episode slices.
type Episodes []Episode

// Last returns the last aired episode from the given slice of Episodes and true
// if it was possible to find such episode.
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

// Next returns the next episode to air from the given slice of Episodes and true
// if it was possible to find such episode.
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

// resultSeason is an internal intermediate struct used for processing EpisodeList results.
type resultSeason struct {
	Number   int       `xml:"no,attr"`
	Episodes []Episode `xml:"episode"`
}

// resultEpisodeList is an internal final struct used for processing EpisodeList results.
type resultEpisodeList struct {
	Total   int            `xml:"totalseasons"`
	Seasons []resultSeason `xml:"Episodelist>Season"`
}

// resultSearch is an internal final struct used for processing Search results.
type resultSearch struct {
	Shows []Show `xml:"show"`
}

const (
	SEARCHURL = `http://services.tvrage.com/feeds/full_search.php?show=%s` // URL for show searching
	EPLISTURL = `http://services.tvrage.com/feeds/episode_list.php?sid=%d` // URL for episode list
	SHOWURL   = `http://services.tvrage.com/feeds/showinfo.php?sid=%d`     // URL for episode list
	TIMEFMT   = `2006-01-02`                                               // time.Parse format string for air date
	VERSION   = `0.0.2`                                                    // library version
)

var (
	Client = &http.Client{} // default HTTP client
)

// parseSearchResult parses the XML as retrieved by Search.
func parseSearchResult(in io.Reader) ([]Show, error) {
	r := resultSearch{}
	x := xml.NewDecoder(in)
	if err := x.Decode(&r); err != nil {
		return nil, err
	}
	return r.Shows, nil
}

func parseGetResult(in io.Reader) (*Show, error) {
	r := Show{}
	x := xml.NewDecoder(in)
	if err := x.Decode(&r); err != nil {
		return nil, err
	}
	r.Name = r.ShowName
	r.Link = r.ShowLink
	r.Country = r.OriginCountry
	return &r, nil
}

// Search retrieves matched shows for the given name.
func Search(name string) ([]Show, error) {
	q := fmt.Sprintf(SEARCHURL, url.QueryEscape(name))
	r, err := Client.Get(q)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return parseSearchResult(r.Body)
}

// Get retrieves the show with the given id
func Get(id int64) (*Show, error) {
	q := fmt.Sprintf(SHOWURL, id)
	r, err := Client.Get(q)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return parseGetResult(r.Body)
}

// parseEpisodeListResults parses the XML as retrieved by EpisodeList.
// It fills in the season number and returns a slice of Episodes.
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

// EpisodeList retrieves the list of episodes for the given show id.
func EpisodeList(id int) (Episodes, error) {
	q := fmt.Sprintf(EPLISTURL, id)
	r, err := Client.Get(q)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return parseEpisodeListResult(r.Body)
}
