package eztvag

import (
	"errors"
	_ "fmt"
	"github.com/arbovm/levenshtein"
	_ "github.com/kr/pretty"
	"gopkg.in/xmlpath.v2"
	"net/http"
	"strings"
)

const DefaultEndpoint = "https://eztv.ag"

var (
	ErrMissingTvShow     = errors.New("eztvag: missing tv show")
	ErrUnexpectedContent = errors.New("eztvag: unexpected content")
	ErrNetworkRequest    = errors.New("eztvag: remote server error")
)

type Torrent struct {
	Name       string
	TorrentURL string
	MagnetURL  string
	Size       string
	Age        string
}

type Tvshow struct {
	Name   string
	Status string
	Rating string
	URL    string
	ImdbId string
}

// Client represents the kickass client
type Client struct {
	Endpoint   string
	HTTPClient *http.Client
	Tvshows    []*Tvshow
}

// New creates a new client
func New() *Client {
	return &Client{
		Endpoint:   DefaultEndpoint,
		HTTPClient: http.DefaultClient,
	}
}

func (c *Client) Init() error {
	return c.loadTvShows()
}

func (c *Client) loadTvShows() error {
	URL := c.Endpoint + "/showlist/"
	resp, err := c.HTTPClient.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ErrNetworkRequest
	}

	root, err := xmlpath.ParseHTML(resp.Body)
	if err != nil {
		return err
	}
	tvshows, err := parseResultShowlist(root)
	if err != nil {
		return err
	}
	if len(tvshows) == 1 {
		return ErrUnexpectedContent
	}
	c.Tvshows = tvshows
	return nil
}

func (c *Client) GetTvShow(name string) ([]*Torrent, error) {
	var guessTvshow *Tvshow
	guessDist := 1000
	for _, v := range c.Tvshows {
		dist := levenshtein.Distance(strings.ToLower(v.Name), strings.ToLower(name))
		if dist < guessDist {
			guessTvshow = v
			guessDist = dist
		}
		if dist == 0 {
			// Perfect match we can stop
			break
		}
	}

	URL := c.Endpoint + guessTvshow.URL
	resp, err := c.HTTPClient.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, ErrNetworkRequest
	}
	root, err := xmlpath.ParseHTML(resp.Body)
	if err != nil {
		return nil, err
	}
	torrents, imdbId, err := parseResultShow(root)
	if err != nil {
		return nil, err
	}
	guessTvshow.ImdbId = imdbId
	return torrents, nil
}
