package eztvag

import (
	_ "fmt"
	"gopkg.in/xmlpath.v2"
	"strings"
)

var (
	xpathShowlistResults = xmlpath.MustCompile("//tr[@name=\"hover\"]")
	xpathShowlistName    = xmlpath.MustCompile(".//a[@class=\"thread_link\"]")
	xpathShowlistStatus  = xmlpath.MustCompile(".//td[2]")
	xpathShowlistRating  = xmlpath.MustCompile(".//td[3]")
	xpathShowlistURL     = xmlpath.MustCompile(".//a[@class=\"thread_link\"]/@href")

	xpathShowResults    = xmlpath.MustCompile("//tr[@name=\"hover\"]")
	xpathShowName       = xmlpath.MustCompile(".//a[@class=\"epinfo\"]")
	xpathShowMagnetURL  = xmlpath.MustCompile(".//a[@class=\"magnet\"]/@href")
	xpathShowTorrentURL = xmlpath.MustCompile(".//a[contains(@class,\"download_\")]/@href")
	xpathShowAge        = xmlpath.MustCompile(".//td[5]")
	xpathShowSize       = xmlpath.MustCompile(".//td[4]")
)

func parseResultShow(root *xmlpath.Node) ([]*Torrent, error) {
	torrents := []*Torrent{}

	iter := xpathShowResults.Iter(root)
	for iter.Next() {
		name, ok := xpathShowName.String(iter.Node())
		if !ok {
			return nil, ErrUnexpectedContent
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
			return nil, ErrUnexpectedContent
		}
		size, ok := xpathShowSize.String(iter.Node())
		if !ok {
			return nil, ErrUnexpectedContent
		}
		t := &Torrent{
			Name:       name,
			MagnetURL:  magnetURL,
			TorrentURL: torrentURL,
			Age:        age,
			Size:       size,
		}
		torrents = append(torrents, t)
	}
	return torrents, nil
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
