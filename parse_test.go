package main

import (
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func TestParseURL(t *testing.T) {
	str := "https://www.youtube.com/embed/aKqPpRQA9pY"
	id := str[strings.LastIndex(str, "/"):]
	if id != "aKqPpRQA9pY" {
		t.Error("parse mismatch", id, "aKqPpRQA9pY")
	}
}

func TestParseList(t *testing.T) {
	t.Log("start list parsing..")

	file, err := ioutil.ReadFile("temp/list.html")
	if err != nil {
		t.Error(err)
	}
	// if err != nil { ... }

	now := time.Now()
	urlList, _, err := parseList(file, time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()))
	if err != nil {
		t.Error(err)
	}
	t.Log(urlList)
}

func TestParseContents(t *testing.T) {
	t.Log("start contents parsing..")

	file, err := ioutil.ReadFile("temp/contents.html")
	if err != nil {
		t.Error(err)
	}
	// if err != nil { ... }

	urlList, err := parseContents(file)
	if err != nil {
		t.Error(err)
	}
	t.Log(urlList)
}
