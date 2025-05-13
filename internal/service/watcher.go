package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FolderWatcher struct {
	folder string
	syncer *DifySyncer
}

func NewFolderWatcher(folder string, syncer *DifySyncer) *FolderWatcher {
	return &FolderWatcher{
		folder: folder,
		syncer: syncer,
	}
}

func (w *FolderWatcher) SyncFolder() error {
	// Get all datasets
	datasets, err := w.syncer.GetDatasets()
	if err != nil {
		return fmt.Errorf("failed to get datasets: %v", err)
	}

	// Create a map of dataset names to IDs
	datasetMap := make(map[string]string)
	for _, dataset := range datasets {
		datasetMap[dataset.Name] = dataset.ID
	}

	// Walk through the folder
	return filepath.Walk(w.folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root folder
		if path == w.folder {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(w.folder, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Handle directory
			return w.handleDirectory(relPath, datasetMap)
		} else {
			// Handle file
			return w.handleFile(path, relPath, datasetMap)
		}
	})
}

func (w *FolderWatcher) handleDirectory(relPath string, datasetMap map[string]string) error {
	// Check if dataset exists
	datasetID, exists := datasetMap[relPath]
	if !exists {
		// Create new dataset
		var err error
		datasetID, err = w.syncer.CreateDataset(relPath)
		if err != nil {
			return fmt.Errorf("failed to create dataset for directory %s: %v", relPath, err)
		}
		datasetMap[relPath] = datasetID
	}
	return nil
}

func (w *FolderWatcher) handleFile(path, relPath string, datasetMap map[string]string) error {
	// Get directory path
	dir := filepath.Dir(relPath)
	if dir == "." {
		return nil // Skip files in root directory
	}

	// Get dataset ID
	datasetID, exists := datasetMap[dir]
	if !exists {
		return fmt.Errorf("parent directory %s not found in datasets", dir)
	}

	// Read file content
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", path, err)
	}

	// Get documents in the dataset
	documents, err := w.syncer.GetDocuments(datasetID)
	if err != nil {
		return fmt.Errorf("failed to get documents for dataset %s: %v", datasetID, err)
	}

	// Check if document exists
	fileName := filepath.Base(path)
	var documentID string
	for _, doc := range documents {
		if doc.Name == fileName {
			documentID = doc.ID
			break
		}
	}

	if documentID != "" {
		// Update existing document
		err = w.syncer.UpdateDocument(datasetID, documentID, fileName, string(content))
		if err != nil {
			return fmt.Errorf("failed to update document %s: %v", fileName, err)
		}
	} else {
		// Create new document
		_, err = w.syncer.CreateDocument(datasetID, fileName, string(content))
		if err != nil {
			return fmt.Errorf("failed to create document for file %s: %v", fileName, err)
		}
	}

	return nil
} 