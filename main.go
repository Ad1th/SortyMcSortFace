package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var extensionGroups = map[string][]string{
	"images": {
		".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".svg", ".heic", ".tiff", ".ico", ".avif",
	},
	"videos": {
		".mp4", ".mkv", ".avi", ".mov", ".wmv", ".webm", ".flv", ".m4v", ".mpeg", ".mpg", ".3gp",
	},
	"music": {
		".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a", ".aiff", ".wma", ".alac",
	},
	"docs": {
		".pdf", ".doc", ".docx", ".txt", ".rtf", ".md", ".odt", ".epub", ".tex", ".log",
	},
	"presentations": {
		".ppt", ".pptx", ".key", ".odp",
	},
	"spreadsheets": {
		".xls", ".xlsx", ".csv", ".ods", ".tsv",
	},
	"archives": {
		".zip", ".rar", ".7z", ".tar", ".gz", ".bz2", ".xz", ".zst", ".lz", ".lz4",
	},
	"code": {
		".go", ".js", ".ts", ".jsx", ".tsx", ".py", ".java", ".c", ".h", ".hpp", ".cpp", ".cs", ".rs", ".php", ".rb", ".swift", ".kt", ".scala", ".sh", ".bash", ".zsh", ".sql", ".json", ".yaml", ".yml", ".toml", ".xml",
	},
	"installers": {
		".exe", ".msi", ".dmg", ".pkg", ".apk", ".deb", ".rpm", ".appimage",
	},
	"fonts": {
		".ttf", ".otf", ".woff", ".woff2",
	},
}

var compoundExtensionCategory = map[string]string{
	".tar.gz":  "archives",
	".tar.bz2": "archives",
	".tar.xz":  "archives",
	".tar.zst": "archives",
}

var keywordCategoryHints = map[string][]string{
	"finance":      {"invoice", "receipt", "tax", "budget", "statement", "salary", "payroll"},
	"screenshots":  {"screenshot", "screen shot", "screen_recording", "screen recording"},
	"design":       {"mockup", "wireframe", "prototype", "figma", "ui-kit", "styleguide"},
	"certificates": {"certificate", "cert", "diploma", "license"},
	"backups":      {"backup", "bak", "snapshot", "restore", "dump"},
}

type options struct {
	dir       string
	recursive bool
	dryRun    bool
	verbose   bool
}

func buildExtensionIndex() map[string]string {
	index := make(map[string]string)
	for category, exts := range extensionGroups {
		for _, ext := range exts {
			index[ext] = category
		}
	}
	return index
}

func sniffCategory(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, readErr := file.Read(buf)
	if readErr != nil && readErr != io.EOF {
		return ""
	}

	contentType := http.DetectContentType(buf[:n])
	if strings.HasPrefix(contentType, "image/") {
		return "images"
	}
	if strings.HasPrefix(contentType, "video/") {
		return "videos"
	}
	if strings.HasPrefix(contentType, "audio/") {
		return "music"
	}
	if strings.Contains(contentType, "pdf") || strings.HasPrefix(contentType, "text/") {
		return "docs"
	}

	return ""
}

func categoryForFile(path string, extensionIndex map[string]string) string {
	name := filepath.Base(path)
	lowerName := strings.ToLower(name)
	ext := strings.ToLower(filepath.Ext(name))

	for compoundExt, category := range compoundExtensionCategory {
		if strings.HasSuffix(lowerName, compoundExt) {
			return category
		}
	}

	if category, ok := extensionIndex[ext]; ok {
		return category
	}

	base := strings.TrimSuffix(lowerName, ext)
	for category, keywords := range keywordCategoryHints {
		for _, keyword := range keywords {
			if strings.Contains(base, keyword) {
				return category
			}
		}
	}

	if sniffed := sniffCategory(path); sniffed != "" {
		return sniffed
	}

	return "others"
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

func gatherFiles(root string, recursive bool) ([]string, error) {
	if !recursive {
		entries, err := os.ReadDir(root)
		if err != nil {
			return nil, err
		}

		files := make([]string, 0, len(entries))
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			files = append(files, filepath.Join(root, entry.Name()))
		}
		return files, nil
	}

	files := make([]string, 0)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})

	return files, err
}

func parseOptions() options {
	var opts options
	flag.StringVar(&opts.dir, "dir", "sort", "Directory to organize")
	flag.BoolVar(&opts.recursive, "recursive", false, "Recursively process files in subfolders")
	flag.BoolVar(&opts.dryRun, "dry-run", false, "Preview actions without moving files")
	flag.BoolVar(&opts.verbose, "verbose", true, "Print per-file actions")
	flag.Parse()
	return opts
}

func main() {
	opts := parseOptions()

	files, err := gatherFiles(opts.dir, opts.recursive)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	extensionIndex := buildExtensionIndex()
	movedCount := 0
	skippedCount := 0
	errorCount := 0
	categoryCount := make(map[string]int)

	for _, srcPath := range files {
		name := filepath.Base(srcPath)
		targetFolder := categoryForFile(srcPath, extensionIndex)
		destDir := filepath.Join(opts.dir, targetFolder)

		if filepath.Dir(srcPath) == destDir {
			skippedCount++
			continue
		}

		destPath := uniqueDestinationPath(destDir, name)

		if opts.verbose {
			fmt.Println("Plan:", srcPath, "->", destPath)
		}

		if opts.dryRun {
			categoryCount[targetFolder]++
			movedCount++
			continue
		}

		if err := os.MkdirAll(destDir, 0o755); err != nil {
			fmt.Println("Error creating directory:", err)
			errorCount++
			continue
		}

		if err := os.Rename(srcPath, destPath); err != nil {
			fmt.Println("Error moving file:", err)
			errorCount++
			continue
		}

		categoryCount[targetFolder]++
		movedCount++
	}

	if opts.dryRun {
		fmt.Println("Dry run complete. No files were moved.")
	} else {
		fmt.Println("Done organizing!")
	}

	fmt.Printf("Summary: moved=%d skipped=%d errors=%d\n", movedCount, skippedCount, errorCount)
	for category, count := range categoryCount {
		fmt.Printf("  - %s: %d\n", category, count)
	}
}
