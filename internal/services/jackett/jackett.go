package jackett

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"golang.org/x/net/context"
)

var (
	apiURL string
	apiKey string
)

type Settings struct {
	ApiURL string
	ApiKey string
	Client *http.Client
}

type FetchRequest struct {
	Query      string
	Trackers   []string
	Categories []uint
}

type FetchResponse struct {
	Results  []Result
	Indexers []Indexer
}

type jackettTime struct {
	time.Time
}

func (jt *jackettTime) UnmarshalJSON(b []byte) (err error) {
	str := strings.Trim(string(b), `"`)
	if str == "0001-01-01T00:00:00" {
	} else if len(str) == 19 {
		jt.Time, err = time.Parse(time.RFC3339, str+"Z")
	} else {
		jt.Time, err = time.Parse(time.RFC3339, str)
	}
	return
}

type Result struct {
	BannerUrl            string
	BlackholeLink        string
	Category             []uint
	CategoryDesc         string
	Comments             string
	Description          string
	DownloadVolumeFactor float32
	Files                uint
	FirstSeen            jackettTime
	Gain                 float32
	Grabs                uint
	Guid                 string
	Imdb                 uint
	InfoHash             string
	Link                 string
	MagnetUri            string
	MinimumRatio         float32
	MinimumSeedTime      uint
	Peers                uint
	PublishDate          jackettTime
	RageID               uint
	Seeders              uint
	Size                 uint
	TMDb                 uint
	TVDBId               uint
	Title                string
	Tracker              string
	TrackerId            string
	UploadVolumeFactor   float32
}

type Indexer struct {
	Error   string
	ID      string
	Name    string
	Results uint
	Status  uint
}

type Jackett struct {
	settings *Settings
}

type SimpleResult struct {
	Title       string `json:"title"`
	Category    []uint `json:"category"`
	MagnetUri   string `json:"magnetUri"`
	Seeders     uint   `json: "seeders`
	Size        uint   `json: "size"`
	Peers       uint   `json: "peers"`
	Description string `json: "description"`
	Tracker     string `json: "Tracker"`
}

func NewJackett(s *Settings) *Jackett {
	if s.ApiURL == "" && apiURL != "" {
		s.ApiURL = apiURL
	}
	if s.ApiKey == "" && apiKey != "" {
		s.ApiKey = apiKey
	}
	if s.Client == nil {
		s.Client = http.DefaultClient
	}
	return &Jackett{settings: s}
}

func (j *Jackett) generateFetchURL(fr *FetchRequest) (string, error) {
	u, err := url.Parse(j.settings.ApiURL)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse apiURL %q", j.settings.ApiURL)
	}
	u.Path = "/api/v2.0/indexers/all/results"
	q := u.Query()
	q.Set("apikey", j.settings.ApiKey)
	for _, t := range fr.Trackers {
		q.Add("Tracker[]", t)
	}
	for _, c := range fr.Categories {
		q.Add("Category[]", fmt.Sprintf("%v", c))
	}
	if fr.Query != "" {
		q.Add("Query", fr.Query)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (j *Jackett) Fetch(ctx context.Context, fr *FetchRequest) (*FetchResponse, error) {
	u, err := j.generateFetchURL(fr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate fetch url")
	}
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make fetch request")
	}
	res, err := j.settings.Client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to invoke fetch request")
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read fetch data")
	}
	var fres FetchResponse
	err = json.Unmarshal(data, &fres)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal fetch data with url=%v and data=%v", u, string(data))
	}
	return &fres, nil
}

func (j *Jackett) FilterResults(results []Result, safeOnly int) ([]byte, error) {
	simpleResults := make([]SimpleResult, 0, len(results))
	for _, r := range results {
		if safeOnly == 1 && r.Tracker != "Internet Archive" {
			continue
		}
		sr := SimpleResult{
			Title:     r.Title,
			Category:  r.Category,
			MagnetUri: r.MagnetUri,
			// Link:      r.Link,
			Tracker: r.Tracker,
			// TrackerId: r.TrackerId,
			Seeders:     r.Seeders,
			Size:        r.Size,
			Peers:       r.Peers,
			Description: r.Description,
		}
		if r.MagnetUri == "" {
			continue
		}
		simpleResults = append(simpleResults, sr)

	}
	return json.Marshal(simpleResults)
}

func init() {
	if v, ok := os.LookupEnv("JACKETT_API_URL"); ok {
		apiURL = v
	}
	if v, ok := os.LookupEnv("JACKETT_API_KEY"); ok {
		apiKey = v
	}
}
