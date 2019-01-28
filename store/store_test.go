package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewRepository(t *testing.T) {
	path := filepath.Clean("./repository")
	repository := NewRepository(path)
	_, err := os.Stat(filepath.Clean(path))
	if os.IsNotExist(err) {
		t.Error("repository has not been created")
	}
	if repository.outDir != filepath.Clean(path) {
		t.Errorf("Expected %v, got %v for outDir\n", path, repository.outDir)
	}
}

func TestConfigurationFilePath(t *testing.T) {
	path := filepath.Clean("./repository")
	repository := NewRepository(path)
	configurationFilePath := filepath.Join(repository.outDir, "config.json")
	if repository.ConfigurationFilePath() != configurationFilePath {
		t.Errorf("Expected %v, got %v for ConfigurationFilePath\n", configurationFilePath, repository.ConfigurationFilePath())
	}
}

func TestAppendDataPath(t *testing.T) {
	path := filepath.Clean("./repository")
	repository := NewRepository(path)
	dataDir := filepath.Join(repository.outDir, "data")
	if repository.AppendDataPath("") != dataDir {
		t.Errorf("Expected %v, got %v for ConfigurationFilePath\n", dataDir, repository.AppendDataPath(""))
	}
}

func TestLogFile(t *testing.T) {
	path := filepath.Clean("./repository")
	repository := NewRepository(path)
	repository.LogFile()

	if _, err := os.Stat(repository.logFilePath()); os.IsNotExist(err) {
		t.Errorf("failed in creating log file at %v\n", repository.logFilePath())
	}
}
