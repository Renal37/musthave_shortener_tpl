package dump_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/dump"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
)

type ShortCollector struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

func TestFillFromStorage(t *testing.T) {
	// Create a temporary file to simulate the input file
	tempFile, err := ioutil.TempFile("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Prepare test data
	events := []ShortCollector{
		{OriginalURL: "http://example.com", ShortURL: "http://short.url/1"},
		{OriginalURL: "http://example.org", ShortURL: "http://short.url/2"},
	}

	// Write test data to the temporary file
	encoder := json.NewEncoder(tempFile)
	for _, event := range events {
		if err := encoder.Encode(event); err != nil {
			t.Fatalf("Failed to encode event: %v", err)
		}
	}
	tempFile.Close()

	// Create a new storage instance
	storageInstance := storage.NewStorage()

	// Call the function under test
	err = dump.FillFromStorage(storageInstance, tempFile.Name())
	if err != nil {
		t.Fatalf("FillFromStorage returned an error: %v", err)
	}

	// Verify that the storage contains the expected data
	for _, event := range events {
		shortURL, exists := storageInstance.Get(event.OriginalURL)
		if !exists {
			t.Errorf("Expected URL %s to exist in storage", event.OriginalURL)
		}
		if shortURL != event.ShortURL {
			t.Errorf("Expected short URL %s, got %s", event.ShortURL, shortURL)
		}
	}
}
