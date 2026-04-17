package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir := "sort" 

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		ext := strings.ToLower(filepath.Ext(name))

		var targetFolder string

		switch ext {
		case ".jpg", ".jpeg", ".png", ".gif":
			targetFolder = "images"
		case ".mp4", ".mkv", ".avi":
			targetFolder = "videos"
		case ".pdf", ".docx", ".txt":
			targetFolder = "docs"
		case ".mp3", ".wav":
    		targetFolder = "music"
		default:
			targetFolder = "others"
		}

		srcPath := filepath.Join(dir, name)
		destDir := filepath.Join(dir, targetFolder)
		destPath := filepath.Join(destDir, name)

		// create folder if not exists
		os.MkdirAll(destDir, os.ModePerm)

		// move file
		err := os.Rename(srcPath, destPath)
		if err != nil {
			fmt.Println("Error moving file:", err)
			continue
		}

		fmt.Println("Moved:", name, "→", targetFolder)
	}

	fmt.Println("Done organizing!")
}