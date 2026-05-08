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
			err := file.Close()
			if err != nil {
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
		for _, v := range dirEntry {
			info, _ := v.Info()
			if info.IsDir() || printFiles {
				totalVisible++
			}
		}

		printed := 0

		for _, v := range dirEntry {
			fileInfo, _ := v.Info()
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
				fmt.Fprintf(out, "%s%s%s\n", prefix, marker, v.Name())
			} else {
				if fileInfo.Size() == 0 {
					fmt.Fprintf(out, "%s%s%s (empty)\n", prefix, marker, v.Name())
				} else {
					fmt.Fprintf(out, "%s%s%s (%db)\n", prefix, marker, v.Name(), fileInfo.Size())
				}
			}

			if fileInfo.IsDir() {
				dirPath := currPath + "/" + v.Name()
				indent := prefix + "│\t"
				if isLast {
					indent = prefix + "\t"
				}
				err := walk(dirPath, indent)
				if err != nil {
					return err
				}

			}
		}

		return nil
	}

	err := walk(path, "")
	if err != nil {
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
