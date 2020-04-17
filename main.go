package main

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

var site = "https://m.clien.net/service/search/group/clien_all?&sk=%s&sv=%s&po=%d"
var origin = "https://m.clien.net"
var search = []string{"music", "mv", "노래", "음악", "뮤직"}

func main() {
	log.SetOutput(io.MultiWriter(&lumberjack.Logger{
		Filename:   "clien.music.log",
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	}, os.Stdout))

	cur, _ := os.Executable()
	conf, err := loadConfig(filepath.Dir(cur) + "/conf.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("start clien.music", conf)

	if !conf.TestMode {
		if err := InitYoutube(filepath.Dir(cur) + "/client_secret.json"); err != nil {
			log.Fatal(err)
		}

		last, err := GetLast(filepath.Dir(cur) + "/last.json")
		if err != nil {
			log.Println("read last file error", err)
		} else {
			for _, link := range last.Links {
				if err := AddSong(last.PlaylistID, link); err != nil {
					log.Println("add song error", err, link, last.PlaylistID)
				}
			}
			return
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

	if err := login(conf.Auth.UserID, conf.Auth.Passwd); err != nil {
		log.Fatal(err)
	}

	var results []string
	for _, se := range conf.SearchUsers {
		for po := 0; po < limit; po++ {
			u := fmt.Sprintf(site, "id", url.QueryEscape(se), po)
			log.Println("start parse list", u)
			data, err := get(u, nil, nil)
			if err != nil {
				panic(err)
			}

			urls, ignore, err := parseList(data, expire)
			if err != nil {
				panic(err)
			}

			// log.Println(urls, ignore)

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
			data, err := get(u, nil, nil)
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
		var errLinks []string
		for _, link := range results {
			id := link[strings.LastIndex(link, "/")+1:]
			// log.Println("add song", link, id, playlistID)
			if err := AddSong(playlistID, id); err != nil {
				log.Println("add song error", err, link, id, playlistID)
				errLinks = append(errLinks, id)
			}
		}

		if len(errLinks) > 0 {

			if err := AddFile(filepath.Dir(cur)+"/last.json", &Last{
				PlaylistID: playlistID,
				Links:      errLinks,
			}); err != nil {
				log.Println("write last log fail", err)
			}

		}

		log.Printf("https://www.youtube.com/playlist?list=%s", playlistID)
	}

}
