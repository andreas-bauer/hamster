package store

import (
	"os"
	"path/filepath"
	"testing"
)

var path string = filepath.Clean("./repository")
var repository *Repository

func TestMain(m *testing.M) {
	r, err := NewRepository(path)
	if err != nil {
		panic(err)
	} else {
		repository = r
	}

	retCode := m.Run()
	os.RemoveAll(path)
	os.Exit(retCode)
}

func TestNewRepository(t *testing.T) {
	_, err := os.Stat(filepath.Clean(path))
	if os.IsNotExist(err) {
		t.Error("repository has not been created")
	}
	if repository.outDir != filepath.Clean(path) {
		t.Errorf("Expected %v, got %v for outDir\n", path, repository.outDir)
	}
}

func TestConfigurationFilePath(t *testing.T) {
	configurationFilePath := filepath.Join(repository.outDir, "config.json")
	if repository.ConfigurationFilePath() != configurationFilePath {
		t.Errorf("Expected %v, got %v for ConfigurationFilePath\n", configurationFilePath, repository.ConfigurationFilePath())
	}
}
