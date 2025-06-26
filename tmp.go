package main

import (
	"github.com/charmbracelet/log"
	"os"
	"path/filepath"
)

func clearTempDir() {
	tempDir := os.TempDir()
	log.Info("clearing temp folder: " + tempDir)

	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Warn("error accessing: " + path + " - " + err.Error())
			return nil
		}
		if path == tempDir {
			return nil
		}
		if info.IsDir() {
			err = os.RemoveAll(path)
		} else {
			err = os.Remove(path)
		}
		if err != nil {
			log.Warn("failed to delete: " + path + " - " + err.Error())
		}
		return nil
	})
	if err != nil {
		log.Error("failed to clear temp folder: " + err.Error())
	} else {
		log.Info("temp folder cleared.")
	}
}
