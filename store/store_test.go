package store

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var path string = filepath.Clean("./repository")
var repository *Repository
var testBytes []byte = []byte("TESTSTRING")

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

func pathExists(path string) bool {
	_, err := os.Stat(filepath.Clean(path))
	return os.IsNotExist(err)
}

func TestNewRepository(t *testing.T) {
	if pathExists(path) {
		t.Error("repository has not been created")
	}
	if repository.outDir != filepath.Clean(path) {
		t.Errorf("Expected %v, got %v for outDir\n", path, repository.outDir)
	}
}

func TestLogFile(t *testing.T) {
	repository.LogFile()
	_, err := ioutil.ReadFile(filepath.Join(repository.outDir, "crawl.log"))
	if err != nil {
		t.Error("log file does not exist")
	}
}

func eval(expected []byte, path string, t *testing.T) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error("file does not exist")
	} else {
		if !bytes.Equal(b, testBytes) {
			t.Errorf("expected %v, got %v", expected, b)

		}
	}
}

func TestStorePayload(t *testing.T) {

	err := repository.StorePayload("test.json", testBytes)
	if err != nil {
		t.Error(err)
	}

	eval(testBytes, filepath.Join(repository.outDir, "data", "test.json"), t)
}

func TestStoreConfigurationJSON(t *testing.T) {
	err := repository.StoreConfigurationJSON(testBytes)
	if err != nil {
		t.Error(err)
	}
	eval(testBytes, filepath.Join(repository.outDir, "config.json"), t)
}
