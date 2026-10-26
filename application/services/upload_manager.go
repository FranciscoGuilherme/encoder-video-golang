package services

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type VideoUpload struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

func (vu *VideoUpload) UploadObject(objectPath string, client *storage.Client, ctx context.Context) error {
	path := strings.Split(objectPath, os.Getenv("localStoragePath")+"/")
	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}
	defer f.Close()

	wc := client.Bucket(vu.OutputBucket).Object(path[1]).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func (vu *VideoUpload) loadPaths() error {
	err := filepath.Walk(vu.VideoPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			vu.Paths = append(vu.Paths, path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func getClientUpload() (*storage.Client, context.Context, error) {
	var err error
	var client *storage.Client

	ctx := context.Background()
	if os.Getenv("STORAGE_EMULATOR_HOST") != "" {
		client, err = storage.NewClient(ctx, option.WithoutAuthentication(), storage.WithJSONReads())
		if err != nil {
			return nil, nil, err
		}
	} else {
		client, err = storage.NewClient(ctx)
		if err != nil {
			return nil, nil, err
		}
	}

	return client, ctx, nil
}
