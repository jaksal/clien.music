package main

import (
	"bytes"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func parseList(r []byte, expire time.Time) ([]string, bool, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(r))
	if err != nil {
		return nil, false, err
	}

	ignore := true

	var result []string
	doc.Find(".list_item").Each(func(i int, s *goquery.Selection) {
		ignore = false
		t := strings.TrimSpace(s.Find(".list_time").Text())
		if strings.Contains(t, ":") {
			arr := strings.Split(t, ":")
			h, _ := strconv.Atoi(arr[0])
			m, _ := strconv.Atoi(arr[1])
			// today
			now := time.Now()
			n := time.Date(now.Year(), now.Month(), now.Day(), h, m, 0, 0, now.Location())
			if !expire.Before(n) {
				//log.Println("ignore expire date", t, n, expire)
				ignore = true
			}
		} else {
			// yester day
			arr := strings.Split(t, "-")
			m, _ := strconv.Atoi(arr[0])
			d, _ := strconv.Atoi(arr[1])
			// today
			now := time.Now()
			n := time.Date(now.Year(), time.Month(m), d, 0, 0, 1, 0, now.Location())
			if !expire.Before(n) {
				//log.Println("ignore expire date", t, n, expire)
				ignore = true
			}
		}

		if ignore {
			return
		}

		/*
			if link, exist := s.Find(".list_subject").Attr("href"); exist && !strings.Contains(link, "rule") {
				// log.Println("add content list", t, link)
				result = append(result, link)
			}
		*/
		if subject := s.Find(".list_subject"); subject != nil {
			if strings.Contains(subject.Text(), "클래식") {
				return
			}
			if link, exist := subject.Attr("href"); exist && !strings.Contains(link, "rule") {
				log.Println("parse content ...", t, strings.TrimSpace(subject.Text()))

				cdata, err := getHTML(origin + link)
				if err != nil {
					panic(err)
				}
				links, err := parseContents(cdata)
				if err != nil {
					panic(err)
				}
				result = append(result, links...)

				//result = append(result, link)
			}

		}
	})
	return result, ignore, nil
}

func parseContents(r []byte) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(r))
	if err != nil {
		return nil, err
	}

	var result []string
	doc.Find(".video > iframe").Each(func(i int, s *goquery.Selection) {
		if link, exist := s.Attr("src"); exist {
			log.Println("find contents ...", link)
			temps := strings.Split(link, "?")
			result = append(result, temps[0])
		}
	})
	return result, nil
}
