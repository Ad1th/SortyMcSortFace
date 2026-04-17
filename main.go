package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func categoryForFile(name string) string {
	lowerName := strings.ToLower(name)
	ext := strings.ToLower(filepath.Ext(name))

	// Handle multi-part archive extensions before using filepath.Ext.
	if strings.HasSuffix(lowerName, ".tar.gz") || strings.HasSuffix(lowerName, ".tar.bz2") || strings.HasSuffix(lowerName, ".tar.xz") {
		return "archives"
	}

	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".svg", ".heic":
		return "images"
	case ".mp4", ".mkv", ".avi", ".mov", ".wmv", ".webm", ".flv", ".m4v":
		return "videos"
	case ".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a":
		return "music"
	case ".pdf", ".doc", ".docx", ".txt", ".rtf", ".md":
		return "docs"
	case ".ppt", ".pptx", ".key":
		return "presentations"
	case ".xls", ".xlsx", ".csv":
		return "spreadsheets"
	case ".zip", ".rar", ".7z", ".tar", ".gz", ".bz2", ".xz":
		return "archives"
	case ".go", ".js", ".ts", ".py", ".java", ".c", ".cpp", ".rs", ".php", ".rb":
		return "code"
	case ".exe", ".msi", ".dmg", ".pkg", ".apk", ".deb":
		return "installers"
	default:
		return "others"
	}
}

func uniqueDestinationPath(destDir, fileName string) string {
	destPath := filepath.Join(destDir, fileName)
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		return destPath
	}

	ext := filepath.Ext(fileName)
	base := strings.TrimSuffix(fileName, ext)

	for i := 1; ; i++ {
		candidate := filepath.Join(destDir, fmt.Sprintf("%s_(%d)%s", base, i, ext))
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
}

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
		targetFolder := categoryForFile(name)

		srcPath := filepath.Join(dir, name)
		destDir := filepath.Join(dir, targetFolder)
		destPath := uniqueDestinationPath(destDir, name)

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