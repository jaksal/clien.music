package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var site = "https://m.clien.net/service/search/group/clien_all?&sk=%s&sv=%s&po=%d"
var origin = "https://m.clien.net"
var search = []string{"music", "mv", "노래", "음악", "뮤직"}

func main() {
	conf, err := loadConfig("conf.json")
	if err != nil {
		log.Fatal(err)
	}

	if !conf.TestMode {
		if err := InitYoutube(); err != nil {
			log.Fatal(err)
		}
	}

	var expire time.Time
	now := time.Now()
	var limit int
	switch conf.Expire {
	case "today":
		expire = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		limit = 3
	case "yesterday":
		yesterday := now.AddDate(0, 0, -1)
		expire = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, now.Location())
		limit = 5
	case "week":
		lastweek := now.AddDate(0, 0, int(now.Weekday())*(-1)+1)
		expire = time.Date(lastweek.Year(), lastweek.Month(), lastweek.Day(), 0, 0, 0, 0, now.Location())
		limit = 10
	case "lastweek":
		lastweek := now.AddDate(0, 0, int(now.Weekday())*(-1)-7+1)
		expire = time.Date(lastweek.Year(), lastweek.Month(), lastweek.Day(), 0, 0, 0, 0, now.Location())
		limit = 10
	case "month":
		expire = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		limit = 30
	default:
		ds := strings.Split(conf.Expire, "-")
		if len(ds) != 2 {
			log.Fatal("invalid expire", expire)
		}
		m, _ := strconv.Atoi(ds[0])
		d, _ := strconv.Atoi(ds[1])
		expire = time.Date(now.Year(), time.Month(m), d, 0, 0, 0, 0, now.Location())
		limit = 20

	}
	log.Println("expire", conf.Expire, expire)

	var results []string
	for _, se := range conf.SearchUsers {
		for po := 0; po < limit; po++ {
			u := fmt.Sprintf(site, "id", url.QueryEscape(se), po)
			log.Println("start parse list", u)
			data, err := getHTML(u)
			if err != nil {
				panic(err)
			}

			urls, ignore, err := parseList(data, expire)
			if err != nil {
				panic(err)
			}

			if ignore {
				break
			}
			results = append(results, urls...)
		}
	}
	for _, se := range conf.SearchTitles {
		for po := 0; po < limit; po++ {
			u := fmt.Sprintf(site, "title", url.QueryEscape(se), po)
			log.Println("start parse list", u)
			data, err := getHTML(u)
			if err != nil {
				panic(err)
			}

			urls, ignore, err := parseList(data, expire)
			if err != nil {
				panic(err)
			}

			if ignore {
				break
			}
			results = append(results, urls...)
		}
	}

	results = removeDuplicate(results)

	if conf.TestMode {
		log.Println(strings.Join(results, "\n"))
	} else {
		// create new playlist
		playlistTitle := fmt.Sprintf("clien.muzic %s~%s", expire.Format("2006-01-02"), time.Now().Format("2006-01-02"))
		playlistID, err := CreatePlaylist(playlistTitle)
		if err != nil {
			log.Fatal(err)
		}

		// add youtube playlist..
		for _, link := range results {
			id := link[strings.LastIndex(link, "/")+1:]
			// log.Println("add song", link, id, playlistID)
			if err := AddSong(playlistID, id); err != nil {
				log.Println("add song error", err, link, id, playlistID)
			}
		}

		log.Printf("https://www.youtube.com/playlist?list=%s", playlistID)
	}

}
