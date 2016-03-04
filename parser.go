package eztvag

import (
	"encoding/hex"
	_ "fmt"
	"github.com/anacrolix/torrent"
	"gopkg.in/xmlpath.v2"
	"regexp"
	"strconv"
	"strings"
)

var (
	xpathShowlistResults = xmlpath.MustCompile("//tr[@name=\"hover\"]")
	xpathShowlistName    = xmlpath.MustCompile(".//a[@class=\"thread_link\"]")
	xpathShowlistStatus  = xmlpath.MustCompile(".//td[2]")
	xpathShowlistRating  = xmlpath.MustCompile(".//td[3]")
	xpathShowlistURL     = xmlpath.MustCompile(".//a[@class=\"thread_link\"]/@href")

	xpathShowLinkImdb   = xmlpath.MustCompile("//a[contains(@href, \"imdb.com/title/tt\")]/@href")
	xpathShowLinkTvMaze = xmlpath.MustCompile("//a[contains(@href, \"tvmaze.com/shows/\")]/@href")
	xpathShowResults    = xmlpath.MustCompile("//tr[@name=\"hover\"]")
	xpathShowName       = xmlpath.MustCompile(".//a[@class=\"epinfo\"]")
	xpathShowLinkEp     = xmlpath.MustCompile(".//a[@class=\"epinfo\"]/@href")
	xpathShowMagnetURL  = xmlpath.MustCompile(".//a[@class=\"magnet\"]/@href")
	xpathShowTorrentURL = xmlpath.MustCompile(".//a[contains(@class,\"download_\")]/@href")
	xpathShowAge        = xmlpath.MustCompile(".//td[5]")
	xpathShowSize       = xmlpath.MustCompile(".//td[4]")

	regexTorrentSeeds          = regexp.MustCompile(`Seeds:.*<.*?>(\d+)</.*?>`)
	regexTorrentPeers          = regexp.MustCompile(`Peers:.*<.*?>(\d+)</.*?>`)
	regexTorrentReleased       = regexp.MustCompile(`Released:.*<.*?>(.+)<.*?>`)
	regexTorrentFile           = regexp.MustCompile(`Torrent File:.*<.*?>(.+)<.*?>`)
	regexTorrentFileFormat     = regexp.MustCompile(`File Format:.*<.*?>(.+)<.*?>`)
	regexTorrentFileResolution = regexp.MustCompile(`Resolution:.*<.*?>(.+)<.*?>`)
	regexTorrentAspectRatio    = regexp.MustCompile(`Aspect Ratio:.*<.*?>(.+)<.*?>`)
	regexTorrentSeason         = regexp.MustCompile(`Season:</b>(.+)<b>`)
	regexTorrentEpisode        = regexp.MustCompile(`Episode:.*<.*?>(.+)\|`)
)

func parseResultTorrent(content string, t *Torrent) error {
	res := regexTorrentSeeds.FindStringSubmatch(content)
	if len(res) != 2 {
		//return fmt.Errorf("Seeds not found")
	} else {
		val, err := strconv.Atoi(res[1])
		if err != nil {
			return err
		}
		t.Seeds = val
	}

	res = regexTorrentPeers.FindStringSubmatch(content)
	if len(res) != 2 {
		// return fmt.Errorf("Peers not found")
	} else {
		val, err := strconv.Atoi(res[1])
		if err != nil {
			return err
		}
		t.Peers = val
	}

	res = regexTorrentReleased.FindStringSubmatch(content)
	if len(res) != 2 {
		//return fmt.Errorf("Released not found")
	} else {
		t.PubDate = strings.Trim(res[1], " ")
	}

	res = regexTorrentFile.FindStringSubmatch(content)
	if len(res) != 2 {
		//return fmt.Errorf("Released not found")
	} else {
		t.Filename = strings.Trim(res[1], " ")
	}
	res = regexTorrentFileFormat.FindStringSubmatch(content)
	if len(res) != 2 {
		//return fmt.Errorf("Released not found")
	} else {
		t.FileFormat = strings.Trim(res[1], " ")
	}
	res = regexTorrentFileResolution.FindStringSubmatch(content)
	if len(res) != 2 {
		//return fmt.Errorf("Released not found")
	} else {
		t.FileResolution = strings.Trim(res[1], " ")
	}
	res = regexTorrentAspectRatio.FindStringSubmatch(content)
	if len(res) != 2 {
		//return fmt.Errorf("Released not found")
	} else {
		t.FileAspectRatio = strings.Trim(res[1], " ")
	}
	res = regexTorrentSeason.FindStringSubmatch(content)
	if len(res) != 2 {
		//return fmt.Errorf("Released not found")
	} else {
		t.Season = strings.Trim(res[1], "&nbsp;")
		t.Season = strings.Trim(t.Season, " ")
	}
	res = regexTorrentEpisode.FindStringSubmatch(content)
	if len(res) != 2 {
		//return fmt.Errorf("Released not found")
	} else {
		t.Episode = strings.Trim(res[1], " ")
	}

	return nil
}

func parseResultShow(root *xmlpath.Node) ([]*Torrent, string, string, error) {
	torrents := []*Torrent{}

	linkImdb, ok := xpathShowLinkImdb.String(root)
	var imdbId string
	if ok {
		linkImdb = strings.TrimRight(linkImdb, "/")
		pos := strings.LastIndex(linkImdb, "/") + 1
		imdbId = linkImdb[pos:]
	}

	linkMaze, ok := xpathShowLinkTvMaze.String(root)
	var tvmazeId string
	if ok {
		pos := strings.LastIndex(linkMaze, "/")
		linkMaze = linkMaze[:pos]
		pos = strings.LastIndex(linkMaze, "/") + 1
		tvmazeId = linkMaze[pos:]
	}

	iter := xpathShowResults.Iter(root)
	for iter.Next() {
		name, ok := xpathShowName.String(iter.Node())
		if !ok {
			return nil, imdbId, tvmazeId, ErrUnexpectedContent
		}

		link, ok := xpathShowLinkEp.String(iter.Node())
		if !ok {
			return nil, imdbId, tvmazeId, ErrUnexpectedContent
		}

		magnetURL, ok := xpathShowMagnetURL.String(iter.Node())
		torrentURL, ok := xpathShowTorrentURL.String(iter.Node())

		// We should have a torrentURL or/and magnetURL
		// otherwise we skip the row
		if len(magnetURL) == 0 && len(torrentURL) == 0 {
			continue
		}

		age, ok := xpathShowAge.String(iter.Node())
		if !ok {
			return nil, imdbId, tvmazeId, ErrUnexpectedContent
		}
		size, ok := xpathShowSize.String(iter.Node())
		if !ok {
			return nil, imdbId, tvmazeId, ErrUnexpectedContent
		}
		sizeb, err := ParseSize([]byte(size))
		if err != nil {
			// trying some cleaning
			pos := strings.LastIndex(size, "(")
			sizeb, err = ParseSize([]byte(size[pos+1:]))
			if err != nil {
				return nil, imdbId, tvmazeId, err
			}
		}
		var infoHash string
		if magnetURL != "" {
			m, err := torrent.ParseMagnetURI(magnetURL)
			if err != nil {
				return nil, imdbId, tvmazeId, err
			}
			infoHash = hex.EncodeToString(m.InfoHash[:])
		}

		t := &Torrent{
			Name:         name,
			InfoHash:     infoHash,
			MagnetURL:    magnetURL,
			TorrentURL:   torrentURL,
			Age:          age,
			Size:         int64(sizeb),
			ShowImdbId:   imdbId,
			ShowTvmazeId: tvmazeId,
			EpisodeURL:   link,
		}
		torrents = append(torrents, t)
	}
	return torrents, imdbId, tvmazeId, nil
}

func parseResultShowlist(root *xmlpath.Node) ([]*Tvshow, error) {
	tvshows := []*Tvshow{}

	iter := xpathShowlistResults.Iter(root)
	for iter.Next() {
		name, ok := xpathShowlistName.String(iter.Node())
		if !ok {
			return nil, ErrUnexpectedContent
		}
		status, ok := xpathShowlistStatus.String(iter.Node())
		if !ok {
			return nil, ErrUnexpectedContent
		}
		rating, ok := xpathShowlistRating.String(iter.Node())
		if !ok {
			return nil, ErrUnexpectedContent
		}
		rating = strings.Trim(rating, "\n ")
		URL, ok := xpathShowlistURL.String(iter.Node())
		if !ok {
			return nil, ErrUnexpectedContent
		}
		t := &Tvshow{
			Name:   name,
			Status: status,
			Rating: rating,
			URL:    URL,
		}
		tvshows = append(tvshows, t)
	}
	return tvshows, nil
}
