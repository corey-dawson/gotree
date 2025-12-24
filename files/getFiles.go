package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetStructure(path string, showAll bool, maxDepth int, crntLevel int) []string {
	var files []string

	if maxDepth == -1 || crntLevel <= maxDepth {
		entries, err := os.ReadDir(path)
		if err != nil {
			fmt.Printf("Error reading directory %s: %v\n", path, err)
			return files
		}
		if len(entries) == 0 {
			files = append(files, path+string(filepath.Separator))
		}
		for _, entry := range entries {

			if !showAll && strings.HasPrefix(entry.Name(), ".") {
				continue
			}
			if !entry.IsDir() {
				files = append(files, filepath.Join(path, entry.Name()))
			} else if maxDepth != -1 && entry.IsDir() && maxDepth == crntLevel {
				files = append(files, filepath.Join(path, entry.Name())+string(filepath.Separator))
			} else {
				subFiles := GetStructure(filepath.Join(path, entry.Name()), showAll, maxDepth, crntLevel+1)
				files = append(files, subFiles...)
			}
		}
	}
	return files
}
