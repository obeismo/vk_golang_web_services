package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	var walk func(string, string) error
	walk = func(currPath, prefix string) error {
		file, err := os.Open(currPath)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			if err := file.Close(); err != nil {
				fmt.Println("problem closing file")
			}
		}(file)

		dirEntry, err := file.ReadDir(-1)
		if err != nil {
			return err
		}

		sort.Slice(dirEntry, func(i, j int) bool {
			return dirEntry[i].Name() < dirEntry[j].Name()
		})

		totalVisible := 0
		for _, entry := range dirEntry {
			info, _ := entry.Info()
			if info.IsDir() || printFiles {
				totalVisible++
			}
		}

		printed := 0

		for _, entry := range dirEntry {
			fileInfo, err := entry.Info()
			if err != nil {
				return err
			}

			if !fileInfo.IsDir() && !printFiles {
				continue
			}

			printed++
			isLast := printed == totalVisible

			marker := "├───"
			if isLast {
				marker = "└───"
			}

			if fileInfo.IsDir() {
				if _, err := fmt.Fprintf(out, "%s%s%s\n", prefix, marker, entry.Name()); err != nil {
					return err
				}
			} else {
				sizeStr := fmt.Sprintf(" (%db)", fileInfo.Size())
				if fileInfo.Size() == 0 {
					sizeStr = " (empty)"
				}

				if _, err := fmt.Fprintf(out, "%s%s%s%s\n", prefix, marker, entry.Name(), sizeStr); err != nil {
					return err
				}
			}

			if fileInfo.IsDir() {
				dirPath := currPath + "/" + entry.Name()
				indent := prefix + "│\t"
				if isLast {
					indent = prefix + "\t"
				}
				if err := walk(dirPath, indent); err != nil {
					return err
				}
			}
		}

		return nil
	}

	if err := walk(path, ""); err != nil {
		return err
	}

	return nil
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
