package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

var service *youtube.Service

//InitYoutube init youtube .
func InitYoutube() error {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		return fmt.Errorf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeScope)
	if err != nil {
		return fmt.Errorf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)
	service, err = youtube.New(client)
	if err != nil {
		return fmt.Errorf("Create Youtube Service fail %v", err)
	}

	return nil
}

//AddSong add new song
func AddSong(playlistID, songID string) error {
	if service == nil {
		return errors.New("not initialize youtube service")
	}
	call := service.PlaylistItems.Insert("snippet", &youtube.PlaylistItem{
		Snippet: &youtube.PlaylistItemSnippet{
			PlaylistId: playlistID,
			ResourceId: &youtube.ResourceId{
				Kind:    "youtube#video",
				VideoId: songID,
			},
		},
	})

	response, err := call.Do()
	if err != nil {
		return err
	}
	log.Println("add song result", response.Snippet.Title)
	return nil
}

// CreatePlaylist create new playlist
func CreatePlaylist(title string) (string, error) {
	if service == nil {
		return "", errors.New("not initialize youtube service")
	}

	call := service.Playlists.Insert("snippet,status", &youtube.Playlist{
		Snippet: &youtube.PlaylistSnippet{
			Title: title,
		},
		Status: &youtube.PlaylistStatus{
			PrivacyStatus: "public",
		},
	})

	response, err := call.Do()
	if err != nil {
		// /fmt.Println("create playlist error", err)
		return "", err
	}

	log.Println("create new playlist", title, response.Id)
	return response.Id, nil

}
