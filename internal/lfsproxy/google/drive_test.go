package google_test

import (
	"fmt"
	"testing"

	"github.com/GareArc/git-lfs-proxy/internal/lfsproxy/google"
	gdrive "google.golang.org/api/drive/v3"
	"github.com/stretchr/testify/assert"
)

func TestSearchForFile(t *testing.T) {
	tests := []struct {
		name     string
		files    []*gdrive.File
		filename string
		want     *gdrive.File
		wantErr  bool
	}{
		{
			name: "file found in root",
			files: []*gdrive.File{
				{Id: "1", Name: "test.txt"},
				{Id: "2", Name: "other.txt"},
			},
			filename: "test.txt",
			want:     &gdrive.File{Id: "1", Name: "test.txt"},
			wantErr:  false,
		},
		{
			name: "file found in subfolder",
			files: []*gdrive.File{
				{
					Id:       "1",
					Name:     "folder",
					MimeType: "application/vnd.google-apps.folder",
				},
			},
			filename: "test.txt",
			want:     &gdrive.File{Id: "2", Name: "test.txt"},
			wantErr:  false,
		},
		{
			name: "file not found",
			files: []*gdrive.File{
				{Id: "1", Name: "other.txt"},
			},
			filename: "test.txt",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "empty file list",
			files:    []*gdrive.File{},
			filename: "test.txt",
			want:     nil,
			wantErr:  true,
		},
		{
			name: "multiple nested folders",
			files: []*gdrive.File{
				{
					Id:       "1",
					Name:     "folder1",
					MimeType: "application/vnd.google-apps.folder",
				},
				{
					Id:       "2",
					Name:     "folder2",
					MimeType: "application/vnd.google-apps.folder",
				},
			},
			filename: "test.txt",
			want:     &gdrive.File{Id: "2", Name: "test.txt"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockGoogleDriveClient{
				files: tt.files,
				subFiles: []*gdrive.File{
					{Id: "2", Name: "test.txt"},
					{Id: "3", Name: "test.txt"},
				},
			}

			got, err := mockClient.SearchForFile(tt.filename)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Id, got.Id)
				assert.Equal(t, tt.want.Name, got.Name)
			}
		})
	}
}

type mockGoogleDriveClient struct {
	google.GoogleDriveClient
	files    []*gdrive.File
	subFiles []*gdrive.File
}

func (m *mockGoogleDriveClient) ListFilesInWorkDir() ([]*gdrive.File, error) {
	return m.files, nil
}

func (m *mockGoogleDriveClient) ListFilesInFolder(folderId string) ([]*gdrive.File, error) {
	return m.subFiles, nil
}

func (m *mockGoogleDriveClient) SearchForFile(filename string) (*gdrive.File, error) {
	files, _ := m.ListFilesInWorkDir()
	return m.searchForFileHelper(filename, files)
}

func (m *mockGoogleDriveClient) searchForFileHelper(filename string, files []*gdrive.File) (*gdrive.File, error) {
	for _, file := range files {
		if file.Name == filename {
			return file, nil
		}

		if file.MimeType == "application/vnd.google-apps.folder" {
			subFiles, _ := m.ListFilesInFolder(file.Id)
			for _, subFile := range subFiles {
				if subFile.Name == filename {
					return subFile, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("file not found")
}
