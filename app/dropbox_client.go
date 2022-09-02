package zettel_bot

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sort"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

const TOKEN = "sl.BOhfAFWhWpKOiqnO4V23_dUxt8dRX5YJL91TgxAEjJq6w15kdjRCNzw7BnHk0tCwgtex6mBQxRSH0wNXuWjqqkDErc2CietaHhXC56GlEQKr79zbCymFTXzueUZ8g_avGUHwIRI"

func GetZeroLinks() string{
	config := dropbox.Config{
		Token: TOKEN,
		LogLevel: dropbox.LogInfo,
	}
	dbx := files.New(config)
	arg := files.NewListFolderArg("/notes/z-core")
	res, err := dbx.ListFolder(arg)
	if err != nil {
		log.Panic(err)
	}
	var file_string []string
	for _, entry := range res.Entries {
		switch f := entry.(type) {
			case *files.FileMetadata:
				file_string = append(file_string, formatFile(f))
		}
	}
	var result []string
	re := regexp.MustCompile(`00.*.md`)
	for _, f := range file_string {
		z_link := re.FindString(f)
		z_link = z_link[:len(z_link)-3]
		result = append(result, z_link)
	}
	sort.Strings(result)
	return strings.Join(result, "\n")
}

func DownloadFile() string{
	config := dropbox.Config{
		Token: TOKEN,
		LogLevel: dropbox.LogInfo,
	}
	dbx := files.New(config)
	arg := files.NewDownloadArg("/notes/base/manual/postgresql/view.md")
	meta, content, err := dbx.Download(arg)
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, meta.Size)
	readed, _ := content.Read(buf)
	if readed == 0 {
		log.Panic("Should not eq 0")
	}
	// links := getZeroLinks(string(buf))
	return string(buf)
}

func getZeroLinks(s string) []string {
	re := regexp.MustCompile(`00-\w+`)
	return re.FindAllString(s, -1)
}

func GetFiles() []string {
	config := dropbox.Config{
		Token: TOKEN,
		LogLevel: dropbox.LogInfo,
	}
	dbx := files.New(config)
	arg := files.NewListFolderArg("/notes/base")
	arg.Recursive=true
	res, err := dbx.ListFolder(arg)
	if err != nil {
		log.Panic(err)
	}
	var result []string
	for _, entry := range res.Entries {
		switch f := entry.(type) {
			case *files.FileMetadata:
				result = append(result, formatFile(f))
			// case *files.FolderMetadata:
			// 	result = append(result, formatFolderMetadata(f))
		}
	}
	result = getMDPath(result)

	dwn_arg := files.NewDownloadArg(result[0])
	meta, content, err := dbx.Download(dwn_arg)
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, meta.Size)
	readed, _ := content.Read(buf)
	if readed == 0 {
		log.Panic("Should not eq 0")
	}
	links := getZeroLinks(string(buf))
	return links
}

func formatFile(e *files.FileMetadata) string {
	return fmt.Sprintf("%s\t\n", e.PathDisplay)
}

func getMDPath(s []string) []string{
	result := make([]string, 0)
	for _, path := range s {
		if idx := strings.LastIndex(path, ".md"); idx != -1 {
			result = append(result, path[:idx+3])
		}
	}
	return result
}

// func formatFolderMetadata(e *files.FolderMetadata) string {
// 	return fmt.Sprintf("%s\t\n", e.PathDisplay)
// }
