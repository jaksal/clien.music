package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

// Config config data
type Config struct {
	YoutubeSecretFile string   `json:"youtube_secret_file"`
	SearchTitles      []string `json:"search_titles"`
	Expire            string   `json:"expire"`
	TestMode          bool     `json:"test_mode"`
}

func loadConfig(file string) (*Config, error) {
	var conf Config

	var search string
	flag.StringVar(&search, "s", "", "search string delemeter is ,")
	flag.StringVar(&conf.YoutubeSecretFile, "secret", "client_secret.json", "google youtube client secret file")
	flag.BoolVar(&conf.TestMode, "t", false, "test mode")
	flag.StringVar(&conf.Expire, "e", "today", "expire date. today,yesterday,week,month")
	flag.Parse()

	if search != "" {
		conf.SearchTitles = strings.Split(search, ",")
	} else {
		// read config file
		dat, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("config file read error %s", err)
		}
		if err := json.Unmarshal(dat, &conf); err != nil {
			return nil, fmt.Errorf("invalid config data %s", err)
		}
	}

	if len(conf.SearchTitles) == 0 || conf.Expire == "" {
		return nil, fmt.Errorf("invalid config %+v", conf)
	}

	log.Printf("load config %+v", conf)
	return &conf, nil
}
