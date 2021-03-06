package eztvag

import (
	_ "fmt"
	"github.com/kr/pretty"
	"gopkg.in/xmlpath.v2"
	"io/ioutil"
	"os"
	"testing"
)

func TestParseResultTorrent2(t *testing.T) {
	f, err := os.Open("test/torrent2.test")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	torrent := &Torrent{}
	parseResultTorrent(string(content), torrent)
	pretty.Println(torrent)
}

func TestParseResultTorrent(t *testing.T) {
	f, err := os.Open("test/torrent.test")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	torrent := &Torrent{}
	parseResultTorrent(string(content), torrent)
	pretty.Println(torrent)
}

func TestParseResultShowList(t *testing.T) {
	f, err := os.Open("test/showlist.test")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	root, err := xmlpath.ParseHTML(f)
	if err != nil {
		t.Error(err)
	}
	tvshows, err := parseResultShowlist(root)
	if err != nil {
		t.Error(err)
	}
	if len(tvshows) != 1593 {
		t.Errorf("expected 1593 tvshows get %d", len(tvshows))
	}

}

// test when magnet is not in page
func TestParseResultShow2(t *testing.T) {
	f, err := os.Open("test/show2.test")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	root, err := xmlpath.ParseHTML(f)
	if err != nil {
		t.Error(err)
	}
	torrents, imdbId, _, err := parseResultShow(root)
	if err != nil {
		t.Error(err)
	}
	if imdbId != "tt0182576" {
		t.Errorf("expected %q get %q", "tt0182576", imdbId)
	}
	if len(torrents) != 100 {
		t.Errorf("expected 28 torrents get %d", len(torrents))
	}
	pretty.Println(torrents)
}

func TestParseResultShow(t *testing.T) {
	f, err := os.Open("test/show.test")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	root, err := xmlpath.ParseHTML(f)
	if err != nil {
		t.Error(err)
	}
	torrents, imdbId, _, err := parseResultShow(root)
	if err != nil {
		t.Error(err)
	}
	if len(torrents) != 28 {
		t.Errorf("expected 28 torrents get %d", len(torrents))
	}
	if imdbId != "tt3032476" {
		t.Errorf("expected %q get %q", "tt3032476", imdbId)
	}

	tr := torrents[0]
	if tr.Name != "Better Call Saul S02E02 HDTV x264-KILLERS" {
		t.Errorf("expected name \"Better Call Saul S02E02 HDTV x264-KILLERS\" get %q",
			tr.Name)
	}
	if tr.Size != 245115125 {
		t.Errorf("excepted size %f get %f", 245115125, tr.Size)
	}
	if tr.Age != "6d 18h" {
		t.Errorf("excepted age \"6d 18h\" get %q", tr.Age)
	}
	if tr.MagnetURL != "magnet:?xt=urn:btih:444a006772d0ebf6f0ad8ef6c4644658d0920f44&dn=Better.Call.Saul.S02E02.HDTV.x264-KILLERS%5Beztv%5D.mp4%5Beztv%5D&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A80&tr=udp%3A%2F%2Fglotorrents.pw%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969" {
		t.Errorf("excepted magnet \"magnet:?xt=urn:btih:444a006772d0ebf6f0ad8ef6c4644658d0920f44&dn=Better.Call.Saul.S02E02.HDTV.x264-KILLERS%5Beztv%5D.mp4%5Beztv%5D&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A80&tr=udp%3A%2F%2Fglotorrents.pw%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969\" get %q", tr.MagnetURL)
	}
}
