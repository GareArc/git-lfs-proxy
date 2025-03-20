package db_test

import (
	"os"
	"testing"

	"github.com/GareArc/git-lfs-proxy/config"
	"github.com/GareArc/git-lfs-proxy/internal/db"
	"github.com/GareArc/git-lfs-proxy/internal/models"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) string {
	tmpDir := t.TempDir()
	config.CONFIG_DIR = tmpDir
	config.Global = &config.Config{
		DBConfig: config.DBConfig{
			Filename: "test.db",
			Username: "testuser",
			Password: "testpass",
		},
	}
	return tmpDir
}

func TestInitDB(t *testing.T) {
	tmpDir := setupTestDB(t)
	defer os.RemoveAll(tmpDir)

	db.InitDB()
	assert.NotNil(t, db.DB)

	var creds []models.OAuthCredential
	result := db.DB.Find(&creds)
	assert.Nil(t, result.Error)
}

func TestInitDBWithInvalidCredentials(t *testing.T) {
	tmpDir := setupTestDB(t)
	defer os.RemoveAll(tmpDir)

	config.Global.DBConfig.Username = ""
	config.Global.DBConfig.Password = ""

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with invalid credentials")
		}
	}()
	db.InitDB()
}

func TestInitDBWithInvalidPath(t *testing.T) {
	config.CONFIG_DIR = "/nonexistent/path"
	config.Global = &config.Config{
		DBConfig: config.DBConfig{
			Filename: "test.db",
			Username: "testuser",
			Password: "testpass",
		},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with invalid path")
		}
	}()
	db.InitDB()
}

func TestDBMigrations(t *testing.T) {
	tmpDir := setupTestDB(t)
	defer os.RemoveAll(tmpDir)

	db.InitDB()

	// Create a test record
	cred := models.OAuthCredential{
		ProxyOwner: "test_owner",
	}
	result := db.DB.Create(&cred)
	assert.Nil(t, result.Error)

	// Verify record was created
	var count int64
	result = db.DB.Model(&models.OAuthCredential{}).Count(&count)
	assert.Nil(t, result.Error)
	assert.Equal(t, int64(1), count)

	// Verify table schema
	columnTypes, err := db.DB.Migrator().ColumnTypes(&models.OAuthCredential{})
	assert.Nil(t, err)

	var columns []string
	for _, column := range columnTypes {
		columns = append(columns, column.Name())
	}

	expectedColumns := []string{"id", "created_at", "updated_at", "deleted_at",
		"token_access_token", "token_token_type", "token_refresh_token",
		"token_expiry", "proxy_owner"}

	for _, expected := range expectedColumns {
		found := false
		for _, actual := range columns {
			if expected == actual {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected column %s not found", expected)
	}
}

func TestDBConnection(t *testing.T) {
	tmpDir := setupTestDB(t)
	defer os.RemoveAll(tmpDir)

	db.InitDB()

	// Test write
	cred := models.OAuthCredential{
		ProxyOwner: "test_owner",
	}
	result := db.DB.Create(&cred)
	assert.Nil(t, result.Error)
	assert.Greater(t, cred.ID, uint(0))

	// Test read
	var readCred models.OAuthCredential
	result = db.DB.First(&readCred, cred.ID)
	assert.Nil(t, result.Error)
	assert.Equal(t, "test_owner", readCred.ProxyOwner)
}
