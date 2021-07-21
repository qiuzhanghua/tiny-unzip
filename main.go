package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	dest, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if len(os.Args) <= 1 {
		fmt.Printf("Usage : %s <filename>\n", os.Args[0])
		os.Exit(1)
	}
	zipFilename := os.Args[1]
	archive, err := zip.OpenReader(zipFilename)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		//if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		//	panic(err)
		//}

		destFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(destFile, fileInArchive); err != nil {
			panic(err)
		}

		destFile.Close()
		fileInArchive.Close()
	}
}
