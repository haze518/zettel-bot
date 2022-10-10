package zettel_bot

import (
	"regexp"
	"sort"
	"time"
	"log"
	"strings"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

type Storage struct {
	Notes []string
	ZeroLinks []string
	LifetimeSecond int
}

func SaveZeroLinks(dropboxClient *DBClient, storage *Storage) {
	for {
		client := *dropboxClient
		if client.IsInitialized() {
			res, err := client.DownloadFolder("/notes/z-core")
			if err != nil {
				log.Panic(err)
			}
			var file_string []string
			for _, entry := range res.Entries {
				switch f := entry.(type) {
					case *files.FileMetadata:
						file_string = append(file_string, f.PathDisplay)
				}
			}
			z_link_name := make([]string, 0, len(file_string))
			re := regexp.MustCompile(`00.*.md`)
			for _, f := range file_string {
				z_link := re.FindString(f)
				if len(z_link) > 0 {
					z_link = z_link[:len(z_link)-3]
					z_link_name = append(z_link_name, z_link)	
				}
			}
			sort.Strings(z_link_name)
			storage.ZeroLinks = z_link_name
			log.Println("Ready to export zero links")
		}
		time.Sleep(time.Second * time.Duration(storage.LifetimeSecond))
	}
}

func SaveNotes(dropboxClient *DBClient, storage *Storage) {
	for {
		client := *dropboxClient
		if client.IsInitialized() {
			storage.Notes = getNotes(client)
			log.Println("Ready to export notes")	
		}
		time.Sleep(time.Second * time.Duration(storage.LifetimeSecond))
	}
}

func getNotes(dropboxClient DBClient) []string {
	paths := make([]string, 0)
	res, err := dropboxClient.DownloadFolder("/notes/base")
	if err != nil {
		log.Panic(err)
	}
	for _, entry := range res.Entries {
		switch f := entry.(type) {
			case *files.FileMetadata:
				paths = append(paths, f.PathDisplay)
		}
	}
	paths = getMDPaths(paths)
	var result []string
	for _, path := range paths {
		dwn_arg := files.NewDownloadArg(path)
		res := dropboxClient.DownloadFile(dwn_arg)
		if len(res) > 0 {
			result = append(result, res)
		}
	}
	return result
}

func getMDPaths(s []string) []string{
	result := make([]string, 0)
	for _, path := range s {
		if idx := strings.LastIndex(path, ".md"); idx != -1 {
			result = append(result, path[:idx+3])
		}
	}
	return result
}