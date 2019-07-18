package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var site = "https://m.clien.net/service/search/group/clien_all?&sk=title&sv=%s&po=%d"
var origin = "https://m.clien.net"
var search = []string{"music", "음악", "노래"}

func main() {
	var expireParam string
	flag.StringVar(&expireParam, "e", "today", "check expire date!. ex=today,yesterday, week, month")
	flag.Parse()

	var expire time.Time
	now := time.Now()
	var limit int
	switch expireParam {
	case "today":
		expire = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		limit = 3
	case "yesterday":
		expire = time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
		limit = 5
	case "week":
		expire = time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday()), 0, 0, 0, 0, now.Location())
		limit = 10
	case "month":
		expire = time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, now.Location())
		limit = 30
	default:
		log.Fatal("invalid expire", expire)
	}
	log.Println("expire", expireParam, expire)

	var results []string
	for _, se := range search {
		for po := 0; po < limit; po++ {
			log.Println("start parse list", fmt.Sprintf(site, url.QueryEscape(se), po))
			data, err := getHTML(fmt.Sprintf(site, se, po))
			if err != nil {
				log.Fatal(err)
			}

			urls, ignore, err := parseList(data, expire)
			if err != nil {
				log.Fatal(err)
			}

			if urls == nil {
				break
			}

			for _, url := range urls {
				cdata, err := getHTML(origin + url)
				if err != nil {
					log.Fatal(err)
				}
				links, err := parseContents(cdata)
				if err != nil {
					log.Fatal(err)
				}
				results = append(results, links...)
			}

			if ignore {
				break
			}
		}
	}

	log.Println(results)

	// add youtube playlist..

}

func getHTML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	// handle the error if there is one
	if err != nil {
		return nil, err
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return html, nil
}
