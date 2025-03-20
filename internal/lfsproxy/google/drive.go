package google

import (
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2"
	gdrive "google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GoogleDriveClient struct {
	client  *gdrive.Service
	workDir string
	ctx     context.Context
}

func NewGoogleDriveClient(tokenSource oauth2.TokenSource, ctx context.Context) (*GoogleDriveClient, error) {
	service, err := gdrive.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, err
	}
	gdriveClient := &GoogleDriveClient{
		client:  service,
		workDir: "appDataFolder",
		ctx:     ctx,
	}

	// verify if client is valid
	_, err = gdriveClient.ListFilesAll()
	if err != nil {
		return nil, err
	}

	return gdriveClient, nil
}

func (g *GoogleDriveClient) SearchForFile(filename string) (*gdrive.File, error) {
	// starting from workDir, search for the file recursively
	files, err := g.ListFilesInWorkDir()
	if err != nil {
		return nil, err
	}

	return g.searchForFileHelper(filename, files)
}

func (g *GoogleDriveClient) searchForFileHelper(filename string, files []*gdrive.File) (*gdrive.File, error) {
	for _, file := range files {
		if file.Name == filename {
			return file, nil
		}

		if file.MimeType == "application/vnd.google-apps.folder" {
			// search in the folder
			subFiles, err := g.ListFilesInFolder(file.Id)
			if err != nil {
				return nil, err
			}

			return g.searchForFileHelper(filename, subFiles)
		}
	}

	return nil, fmt.Errorf("file not found")
}

func (g *GoogleDriveClient) CreateFile(filename string, data io.Reader) (*gdrive.File, error) {
	fileMetadata := &gdrive.File{
		Name:    filename,
		Parents: []string{g.workDir},
	}

	newFile, err := g.client.Files.Create(fileMetadata).
		Media(data).
		Fields("id, name").
		Do()

	if err != nil {
		return nil, err
	}

	return newFile, nil
}

func (g *GoogleDriveClient) ListFilesAll() ([]*gdrive.File, error) {
	return g.listFiles("")
}

func (g *GoogleDriveClient) ListFilesInWorkDir() ([]*gdrive.File, error) {
	return g.listFiles("'" + g.workDir + "' in parents")
}

func (g *GoogleDriveClient) ListFilesInFolder(folderId string) ([]*gdrive.File, error) {
	return g.listFiles("'" + folderId + "' in parents")
}

func (g *GoogleDriveClient) listFiles(filter string) ([]*gdrive.File, error) {
	files, err := g.client.Files.List().
		Fields("files(id, name)").
		Spaces("appDataFolder"). // enforce search in appDataFolder
		Q(filter).
		Q("trashed = false").
		Do()

	if err != nil {
		return nil, err
	}

	return files.Files, nil
}

func (g *GoogleDriveClient) DownloadFile(fileId string) (*gdrive.File, error) {

	file, err := g.client.Files.Get(fileId).Do()
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (g *GoogleDriveClient) DeleteFile(fileId string) error {
	err := g.client.Files.Delete(fileId).Do()
	if err != nil {
		return err
	}

	return nil
}
