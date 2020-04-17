package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Last struct {
	PlaylistID string   `json:"playlist_id"`
	Links      []string `json:"links"`
}

func AddFile(path string, last *Last) error {

	dat, _ := json.MarshalIndent(last, "", " ")

	return ioutil.WriteFile(path, dat, 0644)

}

func GetLast(path string) (*Last, error) {

	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var last Last
	if err := json.Unmarshal(dat, &last); err != nil {
		return nil, err
	}

	os.Remove(path)

	return &last, nil
}
