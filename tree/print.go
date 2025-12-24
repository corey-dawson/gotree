package tree

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// PrintTree converts a flat list of file paths into a tree structure and prints it
func PrintTree(fileList []string, rootPath string) {
	output := GetTreeString(fileList, rootPath)
	fmt.Print(output)
}

// GetTreeString converts a flat list of file paths into a tree structure and returns it as a string
func GetTreeString(fileList []string, rootPath string) string {
	if len(fileList) == 0 {
		return rootPath + "\n"
	}

	// Build a tree structure from the flat list
	tree := buildTree(fileList, rootPath)

	var builder strings.Builder
	builder.WriteString(rootPath)
	builder.WriteString("\n")

	// Build the tree string recursively
	buildTreeString(tree, "", true, &builder)

	return builder.String()
}

// treeNode represents a node in the file tree
type treeNode struct {
	name     string
	isDir    bool
	children map[string]*treeNode
}

// buildTree constructs a tree structure from a flat list of paths
func buildTree(fileList []string, rootPath string) map[string]*treeNode {
	tree := make(map[string]*treeNode)

	// Normalize root path
	rootPath, _ = filepath.Abs(rootPath)

	for _, filePath := range fileList {
		absPath, _ := filepath.Abs(filePath)

		// Get relative path from root
		relPath, err := filepath.Rel(rootPath, absPath)
		if err != nil {
			continue
		}

		// Skip if path is "." (root itself)
		if relPath == "." {
			continue
		}

		// Determine if it's a directory (ends with separator)
		isDir := strings.HasSuffix(filePath, string(filepath.Separator))

		// Split path into components
		parts := strings.Split(relPath, string(filepath.Separator))

		// Build the tree
		current := tree
		for i, part := range parts {
			isLast := i == len(parts)-1

			if node, exists := current[part]; exists {
				// Node already exists, mark as directory if it has children or is the last part and isDir
				if !isLast || isDir {
					node.isDir = true
				}
				if node.children == nil {
					node.children = make(map[string]*treeNode)
				}
				current = node.children
			} else {
				// Create new node
				// Mark as directory if: it's not the last part (intermediate) OR it's the last part and isDir
				nodeIsDir := !isLast || isDir
				node := &treeNode{
					name:     part,
					isDir:    nodeIsDir,
					children: make(map[string]*treeNode),
				}
				current[part] = node

				if !isLast {
					current = node.children
				}
			}
		}
	}

	return tree
}

// printTreeRecursive prints the tree with proper box-drawing characters
func printTreeRecursive(nodes map[string]*treeNode, prefix string, isLast bool) {
	if len(nodes) == 0 {
		return
	}

	keys := getSortedKeys(nodes)

	for i, key := range keys {
		node := nodes[key]
		isLastChild := i == len(keys)-1

		// Determine the connector character
		var connector string
		if isLastChild {
			connector = "└── "
		} else {
			connector = "├── "
		}

		// Print the current node
		fmt.Print(prefix + connector + node.name)
		if node.isDir {
			fmt.Print("/")
		}
		fmt.Println()

		// Prepare prefix for children
		// If this node is the last child, use spaces; otherwise continue the vertical line
		var childPrefix string
		if isLastChild {
			childPrefix = prefix + "    "
		} else {
			childPrefix = prefix + "│   "
		}

		// Recursively print children
		printTreeRecursive(node.children, childPrefix, isLastChild)
	}
}

// buildTreeString builds the tree string recursively
func buildTreeString(nodes map[string]*treeNode, prefix string, isLast bool, builder *strings.Builder) {
	if len(nodes) == 0 {
		return
	}

	keys := getSortedKeys(nodes)

	for i, key := range keys {
		node := nodes[key]
		isLastChild := i == len(keys)-1

		// Determine the connector character
		var connector string
		if isLastChild {
			connector = "└── "
		} else {
			connector = "├── "
		}

		// Build the current node string
		builder.WriteString(prefix)
		builder.WriteString(connector)
		builder.WriteString(node.name)
		if node.isDir {
			builder.WriteString("/")
		}
		builder.WriteString("\n")

		// Prepare prefix for children
		// If this node is the last child, use spaces; otherwise continue the vertical line
		var childPrefix string
		if isLastChild {
			childPrefix = prefix + "    "
		} else {
			childPrefix = prefix + "│   "
		}

		// Recursively build children
		buildTreeString(node.children, childPrefix, isLastChild, builder)
	}
}

// getSortedKeys returns sorted keys with directories first, then files
func getSortedKeys(nodes map[string]*treeNode) []string {
	keys := make([]string, 0, len(nodes))
	dirKeys := make([]string, 0)
	fileKeys := make([]string, 0)

	for k := range nodes {
		// A node is a directory if it's marked as such OR has children
		if nodes[k].isDir || len(nodes[k].children) > 0 {
			dirKeys = append(dirKeys, k)
		} else {
			fileKeys = append(fileKeys, k)
		}
	}

	// Sort directories and files separately
	sort.Strings(dirKeys)
	sort.Strings(fileKeys)

	// Combine: directories first, then files
	keys = append(keys, dirKeys...)
	keys = append(keys, fileKeys...)

	return keys
}
