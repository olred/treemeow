package main

import (
	"fmt"
	// "io"
	"os"
	"path/filepath"
	"strings"
)

func dirTree(out *os.File, path string, printFiles bool) error {
	var levelsInfo []int
	err := getTree(path, out, 0, printFiles, &levelsInfo)
	if err != nil {
		return err
	}
	return nil
}

// levelsInfo 0 - not last directory, 1 - last directory
func getTree(path string, out *os.File, level int, printFiles bool, levelsInfo *[]int ) error {
	entries, err := os.ReadDir(path)
	var indexFiles []int
	if !printFiles {
		for i, e := range entries {
			if e.IsDir() {
				indexFiles = append(indexFiles, i)
			}
		}
	} else {
		for i, _ := range entries {
			indexFiles = append(indexFiles, i)
		}
	}
	if err != nil {
		return err
	}
	tabs := strings.Repeat(" ", 4)
	for i, e := range indexFiles {
		fileName := entries[e].Name()
		filePlace := filepath.Join(path, fileName)

		if err != nil {
			return err
		}
		
		if i == len(indexFiles)-1 {
			if entries[e].IsDir() {
				writer(fileName, levelsInfo, true, "└───", tabs, level)
				getTree(filePlace, out, level + 1, printFiles, levelsInfo)
				continue
			} else if printFiles{
				fileStat, _ := os.Stat(filePlace)
				IntSize := fileStat.Size()
				var fileSize string
				if IntSize == 0 {
					fileSize = "empty"
				} else {
					fileSize = fmt.Sprintf("%db", IntSize)
				}
				writer(fileName, levelsInfo, false, "└───", tabs, level, fileSize)
				continue
			}
		}
		if entries[e].IsDir() {
			writer(fileName, levelsInfo, true, "├───", tabs, level)
			getTree(filePlace, out, level + 1, printFiles, levelsInfo)
		} else if printFiles {
			fileStat, _ := os.Stat(filePlace)
			IntSize := fileStat.Size()
			var fileSize string
			if IntSize == 0 {
				fileSize = "empty"
			} else {
				fileSize = fmt.Sprintf("%db", IntSize)
			}
			writer(fileName, levelsInfo, false, "├───", tabs, level, fileSize)
			continue
		}
	}
	return nil
}

func writer(fileName string, levelsInfo *[]int, isDirectory bool, graph string, tabs string, level int, fileSize ...string) {
	if graph == "├───" {
		if len(*levelsInfo) - 1 >= level {
			(*levelsInfo)[level] = 0
		} else {
			*levelsInfo = append(*levelsInfo, 0)
		}
	} else {
		if len(*levelsInfo) - 1 >= level {
			(*levelsInfo)[level] = 1
		} else {
			*levelsInfo = append(*levelsInfo, 1)
		}
	}
	var linksString string
	tempTabul := 4
	for i := 0; i < level; i++ {
		if (*levelsInfo)[i] == 0 {
			linksString += strings.Repeat(" ", tempTabul) + "^-^"
		} else {
			linksString += strings.Repeat(" ", tempTabul)
		}
	}
	if isDirectory {
		fmt.Printf(linksString + tabs + "%v\x1b[34;1m%v \x1b[0m\n", graph, fileName)

	} else {
		fmt.Printf(linksString + tabs + "%v\x1b[35;1m%v (%v) meow~\x1b[0m\n", graph, fileName, fileSize[0])
	}
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
