package zettel_bot

import (
	"log"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/go-redis/redis"
)

type Storage struct {
	LifetimeSecond int
	RedisClient    *redis.Client
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
			pipe := storage.RedisClient.Pipeline()
			pipe.Del("zero_links")
			for _, val := range z_link_name {
				_, err = pipe.RPush("zero_links", val).Result()
				if err != nil {
					log.Printf("Error occured while save zero_links\n%s", err)
				}
			}
			pipe.Exec()
			log.Println("Ready to export zero links")
		}
		time.Sleep(time.Second * time.Duration(storage.LifetimeSecond))
	}
}

func SaveNotes(dropboxClient *DBClient, storage *Storage) {
	for {
		client := *dropboxClient
		if client.IsInitialized() {
			pipe := storage.RedisClient.Pipeline()
			pipe.Del("notes")
			_, err := pipe.SAdd("notes", getNotes(client)).Result()
			if err != nil {
				log.Printf("Error occured while save notes\n%s", err)
			}
			pipe.Exec()
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

func getMDPaths(s []string) []string {
	result := make([]string, 0)
	for _, path := range s {
		if idx := strings.LastIndex(path, ".md"); idx != -1 {
			result = append(result, path[:idx+3])
		}
	}
	return result
}
