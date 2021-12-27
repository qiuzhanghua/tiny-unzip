package main

import (
	"archive/zip"
	"fmt"
	flag "github.com/spf13/pflag"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var destination string

	flag.Usage = func() {
		fmt.Printf("Unzip File to Destination Folder\n\nUSAGE:\n%s <filename> [OPTIONS]\n\nOPTIONS:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println()
	}
	flag.StringVarP(&destination, "exdir", "d", ".", "Directory where files be extracted into")
	flag.Parse()

	if len(os.Args) <= 1 {
		flag.Usage()
		os.Exit(0)
	}
	if strings.HasSuffix(os.Args[1], "help") {
		flag.Usage()
		os.Exit(0)
	}
	if strings.HasSuffix(os.Args[1], "version") {
		fmt.Printf("tiny-unzip %s (%s %s)\n", AppVersion, AppRevision, AppBuildDate)
		os.Exit(0)
	}

	zipFilename := os.Args[1]
	archive, err := zip.OpenReader(zipFilename)
	if err != nil {
		fmt.Printf("Can't find file named %s!\n", zipFilename)
		os.Exit(1)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(destination, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}
		dir := filepath.Dir(filePath)
		os.MkdirAll(dir, os.ModePerm)

		destFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			fmt.Printf("Open %s error: %s!\n", filePath, err)
			os.Exit(1)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			fmt.Printf("Open %s error: %s!\n", fileInArchive, err)
			os.Exit(1)
		}

		if _, err := io.Copy(destFile, fileInArchive); err != nil {
			fmt.Printf("Copy error: %s!\n", err)
			os.Exit(1)
		}

		destFile.Close()
		fileInArchive.Close()
	}
}
