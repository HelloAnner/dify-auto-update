package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FolderWatcher 文件夹监视器结构体
type FolderWatcher struct {
	folder string
	syncer *DifySyncer
	// 追踪现有文件及其对应的文档ID
	fileDocuments map[string]fileDocument
}

// fileDocument 文件文档关联信息
type fileDocument struct {
	documentID string
	datasetID  string
}

// NewFolderWatcher 创建新的文件夹监视器实例
func NewFolderWatcher(folder string, syncer *DifySyncer) *FolderWatcher {
	return &FolderWatcher{
		folder:        folder,
		syncer:        syncer,
		fileDocuments: make(map[string]fileDocument),
	}
}

// SyncFolder 同步文件夹内容到Dify
func (w *FolderWatcher) SyncFolder() error {
	// 获取所有
	datasets, err := w.syncer.GetDatasets()
	if err != nil {
		return fmt.Errorf("获取数据集失败：%v", err)
	}

	// 创建数据集名称到ID的映射
	datasetMap := make(map[string]string)
	for _, dataset := range datasets {
		datasetMap[dataset.Name] = dataset.ID
	}

	// 追踪当前文件用于检测删除的文件
	currentFiles := make(map[string]bool)

	// 遍历文件夹
	err = filepath.Walk(w.folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过根文件夹
		if path == w.folder {
			return nil
		}

		// 获取相对路径
		relPath, err := filepath.Rel(w.folder, path)
		if err != nil {
			return err
		}

		// 跳过隐藏文件和目录（以.开头的文件和目录）
		if strings.HasPrefix(filepath.Base(path), ".") {
			if info.IsDir() {
				return filepath.SkipDir // 如果是目录，跳过整个目录
			}
			return nil // 如果是文件，跳过这个文件
		}

		if info.IsDir() {
			// 处理目录
			return w.handleDirectory(relPath, datasetMap)
		} else {
			// 标记文件为存在
			currentFiles[relPath] = true
			// 处理文件
			return w.handleFile(path, relPath, datasetMap)
		}
	})

	if err != nil {
		return err
	}

	// 检查已删除的文件
	for filePath, doc := range w.fileDocuments {
		if !currentFiles[filePath] {
			// 文件不再存在，从Dify中删除
			if err := w.syncer.DeleteDocument(doc.datasetID, doc.documentID); err != nil {
				return fmt.Errorf("删除文件 %s 对应的文档失败：%v", filePath, err)
			}
			delete(w.fileDocuments, filePath)
		}
	}

	return nil
}

// handleDirectory 处理目录
func (w *FolderWatcher) handleDirectory(relPath string, datasetMap map[string]string) error {
	// 检查数据集是否存在
	datasetID, exists := datasetMap[relPath]
	if !exists {
		// 创建新的数据集
		var err error
		datasetID, err = w.syncer.CreateDataset(relPath)
		if err != nil {
			return fmt.Errorf("为目录 %s 创建数据集失败：%v", relPath, err)
		}
		datasetMap[relPath] = datasetID
	}
	return nil
}

// handleFile 处理文件
func (w *FolderWatcher) handleFile(path, relPath string, datasetMap map[string]string) error {
	// 获取目录路径
	dir := filepath.Dir(relPath)
	if dir == "." {
		return nil // 跳过根目录下的文件
	}

	// 获取数据集ID
	datasetID, exists := datasetMap[dir]
	if !exists {
		return fmt.Errorf("未找到父目录 %s 对应的数据集", dir)
	}

	// 读取文件内容
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取文件 %s 失败：%v", path, err)
	}

	// 获取数据集中的文档
	documents, err := w.syncer.GetDocuments(datasetID)
	if err != nil {
		return fmt.Errorf("获取数据集 %s 的文档失败：%v", datasetID, err)
	}

	// 检查文档是否存在
	fileName := filepath.Base(path)
	var documentID string
	for _, doc := range documents {
		if doc.Name == fileName {
			documentID = doc.ID
			break
		}
	}

	if documentID != "" {
		// 更新现有文档
		err = w.syncer.UpdateDocument(datasetID, documentID, fileName, string(content))
		if err != nil {
			return fmt.Errorf("更新文档 %s 失败：%v", fileName, err)
		}
		// 更新追踪信息
		w.fileDocuments[relPath] = fileDocument{
			documentID: documentID,
			datasetID:  datasetID,
		}
	} else {
		// 创建新文档
		documentID, err = w.syncer.CreateDocument(datasetID, fileName, string(content))
		if err != nil {
			return fmt.Errorf("为文件 %s 创建文档失败：%v", fileName, err)
		}
		// 追踪新文档
		w.fileDocuments[relPath] = fileDocument{
			documentID: documentID,
			datasetID:  datasetID,
		}
	}

	return nil
} 