package store

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var path string = filepath.Clean("./repository")

func TestMain(m *testing.M) {
	retCode := m.Run()
	os.RemoveAll(path)
	os.Exit(retCode)
}

func createTestRepository(t *testing.T) Repository {
	repository, err := NewRepository(path)
	if err != nil {
		t.Error("repository could not be created")
	}
	return repository
}

func TestNewRepository(t *testing.T) {
	repository := createTestRepository(t)
	_, err := os.Stat(filepath.Clean(path))
	if os.IsNotExist(err) {
		t.Error("repository has not been created")
	}
	if repository.outDir != filepath.Clean(path) {
		t.Errorf("Expected %v, got %v for outDir\n", path, repository.outDir)
	}
}

func TestConfigurationFilePath(t *testing.T) {
	repository := createTestRepository(t)
	configurationFilePath := filepath.Join(repository.outDir, "config.json")
	if repository.ConfigurationFilePath() != configurationFilePath {
		t.Errorf("Expected %v, got %v for ConfigurationFilePath\n", configurationFilePath, repository.ConfigurationFilePath())
	}
}

func TestAppendDataPath(t *testing.T) {
	repository := createTestRepository(t)
	dataDir := filepath.Join(repository.outDir, "data")
	if repository.AppendDataPath("") != dataDir {
		t.Errorf("Expected %v, got %v for ConfigurationFilePath\n", dataDir, repository.AppendDataPath(""))
	}
}

func TestLogFile(t *testing.T) {
	repository := createTestRepository(t)
	repository.LogFile()

	if _, err := os.Stat(repository.logFilePath()); os.IsNotExist(err) {
		t.Errorf("failed in creating log file at %v\n", repository.logFilePath())
	}
}

func TestStore(t *testing.T) {
	repository := createTestRepository(t)
	b := []byte{0, 1, 2, 3, 4}
	writeErr := repository.Store(repository.AppendDataPath("test"), b)
	rb, readErr := ioutil.ReadFile(repository.AppendDataPath("test"))

	if writeErr != nil || readErr != nil {
		t.Error("failed in writing bytes to file")
	}

	if !bytes.Equal(b, rb) {
		t.Errorf("Expected %v, got %v byte array\n", b, rb)
	}
}

func TestFielExists(t *testing.T) {
	repository := createTestRepository(t)
	repository.logFilePath()
	testFilePath := repository.AppendDataPath("test")

	writeErr := repository.Store(testFilePath, []byte{0, 1, 2, 3, 4})

	if writeErr != nil || repository.FileExists(testFilePath) != true {
		t.Error("method FileExists(...) failed")
	}
}
