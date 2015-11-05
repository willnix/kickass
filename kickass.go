package kickass

import (
	"errors"

	"github.com/PuerkitoBio/goquery"
)

// Eztv errors
var (
	ErrMovieNotFound   = errors.New("movie not found")
	ErrEmptyResponse   = errors.New("empty response from server")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrMissingArgument = errors.New("missing argument")
	ErrParsingFailure  = errors.New("parsing error")
)

// Movie contains details and sources of movie
type Movie struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Cover string `json:"cover"`
	// [1080p|720p|hdtv]
	Sources map[string]string `json:"sources"`
}

// GetMovie finds shows with a title containg the keyword
// Returns error if no show is found
func GetMovie(keyword string) (*Movie, error) {
	if keyword == "" {
		return nil, ErrMissingArgument
	}

	doc, err := goquery.NewDocument("https://kat.cr/usearch/" + keyword)
	if err != nil {
		return nil, err
	}

	usearch := doc.Find(".torrentMediaInfo")
	if usearch.Length() < 1 {
		return nil, ErrMovieNotFound
	}

	titleLink := doc.Find("h1 > a.plain")
	title := titleLink.Text()
	if title == "" {
		return nil, ErrParsingFailure
	}

	url, ok := titleLink.Attr("href")
	if !ok {
		return nil, ErrParsingFailure
	}

	doc, err = goquery.NewDocument("https://kat.cr" + url)
	if err != nil {
		return nil, err
	}

	cover, ok := doc.Find(".movieCover > img").Attr("src")
	if !ok {
		return nil, ErrParsingFailure
	}

	magnets := make(map[string]string, 3)
	magnets["1080p"], _ = doc.Find("#tab-1080p i.ka-magnet").Parent().Attr("href")
	magnets["720p"], _ = doc.Find("#tab-720p i.ka-magnet").Parent().Attr("href")
	magnets["hdtv"], _ = doc.Find("#tab-HDRiP i.ka-magnet").Parent().Attr("href")

	return &Movie{
			Title:   title,
			URL:     url,
			Cover:   cover,
			Sources: magnets},
		nil

}
