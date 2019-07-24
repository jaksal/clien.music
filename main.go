package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
)

var site = "https://m.clien.net/service/search/group/clien_all?&sk=title&sv=%s&po=%d"
var origin = "https://m.clien.net"
var search = []string{"music", "mv", "노래", "음악", "뮤직"}

func main() {
	conf, err := loadConfig("conf.json")
	if err != nil {
		log.Fatal(err)
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
		lastweek := now.AddDate(0, 0, int(now.Weekday())*(-1))
		expire = time.Date(lastweek.Year(), lastweek.Month(), lastweek.Day(), 0, 0, 0, 0, now.Location())
		limit = 10
	case "month":
		expire = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		limit = 30
	default:
		log.Fatal("invalid expire", expire)
	}
	log.Println("expire", conf.Expire, expire)

	var results []string
	for _, se := range conf.SearchTitles {
		for po := 0; po < limit; po++ {
			u := fmt.Sprintf(site, url.QueryEscape(se), po)
			log.Println("start parse list", u)
			data, err := getHTML(u)
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

	results = removeDuplicate(results)

	if conf.TestMode {
		log.Println(strings.Join(results, "\n"))
	} else {
		if err := InitYoutube(); err != nil {
			log.Fatal(err)
		}

		// create new playlist
		playlistTitle := fmt.Sprintf("clien.music %s-%s", expire.Format("2006-01-02"), time.Now().Format("2006-01-02"))
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
