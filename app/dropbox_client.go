package zettel_bot

import (
	"log"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)


type DBClient interface {
	DownloadFolder(path string) (*files.ListFolderResult, error)
	DownloadFile(path *files.DownloadArg) string
	IsInitialized() bool
}

type DropboxCLient struct {
	dbx files.Client
}


func NewDropboxClient(token string) *DropboxCLient {
	config := dropbox.Config{
		Token: token,
		LogLevel: dropbox.LogInfo,
	}
	client := files.New(config)
	return &DropboxCLient{client}
}

func (d *DropboxCLient) DownloadFolder(path string) (*files.ListFolderResult, error) {
	arg := files.NewListFolderArg(path)
	arg.Recursive = true
	return d.dbx.ListFolder(arg)
}

func (d *DropboxCLient) DownloadFile(path *files.DownloadArg) string {
	meta, content, err := d.dbx.Download(path)
	if err != nil {
		log.Println(err)
	}
	buf := make([]byte, meta.Size)
	readed, _ := content.Read(buf)
	if readed > 0 {
		return string(buf)
	}
	return ""
}

func (d *DropboxCLient) IsInitialized() bool {
	return d.dbx != nil
}
