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
