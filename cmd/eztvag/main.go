package main

import (
	"github.com/kr/pretty"
	"github.com/ricksancho/eztvag"
	"os"
)

func main() {
	query := os.Args[1]
	e := eztvag.New()
	err := e.Init()
	if err != nil {
		panic(err)
	}
	torrents, err := e.GetTvShow(query)
	if err != nil {
		panic(err)
	}
	pretty.Println(torrents)
}
