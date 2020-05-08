package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

// Config config data
type Config struct {
	YoutubeSecretFile string   `json:"youtube_secret_file"`
	SearchTitles      []string `json:"search_titles"`
	SearchUsers       []string `json:"search_users"`
	Expire            string   `json:"expire"`
	Mode              int      `json:"mode"`
	Auth              struct {
		UserID string `json:"userid"`
		Passwd string `json:"passwd"`
	} `json:"auth"`
}

func loadConfig(file string) (*Config, error) {
	var conf Config

	var mode int
	flag.IntVar(&mode, "m", 0, "test mode 0=real 1=test 2=last")
	flag.Parse()

	// read config file
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("config file read error %s", err)
	}
	if err := json.Unmarshal(dat, &conf); err != nil {
		return nil, fmt.Errorf("invalid config data %s", err)
	}

	if (len(conf.SearchTitles) == 0 && len(conf.SearchUsers) == 0) || conf.Expire == "" {
		return nil, fmt.Errorf("invalid config %+v", conf)
	}

	conf.Mode = mode

	log.Printf("load config %+v", conf)
	return &conf, nil
}
