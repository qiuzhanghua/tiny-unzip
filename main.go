package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/labstack/gommon/log"
	flag "github.com/spf13/pflag"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	log.SetPrefix("unzip")
	format := strings.ToLower(os.Getenv("LOGGING_FORMAT"))
	if format != "json" {
		log.SetHeader(`${prefix}, ${level} ${short_file}(${line})`)
	}
	log.SetOutput(os.Stdout)
	level := strings.ToLower(os.Getenv("LOGGING_LEVEL"))
	x := levelOf(level)
	log.SetLevel(x)
}

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
	err := Unzip(zipFilename, destination)
	if err != nil {
		fmt.Println("Err", err)
	}
}

func Unzip(zipFilename string, destination string) error {
	archive, err := zip.OpenReader(zipFilename)
	if err != nil {
		return err
	}
	defer archive.Close()
	linkMap := make(map[string]string, 0)

	for _, f := range archive.File {
		filePath := filepath.Join(destination, f.Name)
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		dir := filepath.Dir(filePath)
		_ = os.MkdirAll(dir, os.ModePerm)

		fileInArchive, err := f.Open()
		//log.Debugf("%s %s", f.Name, f.FileInfo().Mode())
		if f.Mode()&fs.ModeSymlink > 0 {
			//log.Debug(f.Mode() & fs.ModeSymlink)
			buf := new(bytes.Buffer)
			_, err := io.Copy(buf, fileInArchive)
			if err != nil {
				return err
			}
			linkMap[f.Name] = buf.String()
			continue
		}

		destFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		if _, err := io.Copy(destFile, fileInArchive); err != nil {
			return err
		}

		destFile.Close()
		fileInArchive.Close()
	}
	wd, err := os.Getwd()
	err = os.Chdir(destination)
	if err != nil {
		return err
	}
	for k, v := range linkMap {
		log.Debugf("%s => %s", v, k)
		err = os.Symlink(v, k)
		if err != nil {
			log.Error(err)
		}
	}
	_ = os.Chdir(wd)
	return nil
}

func levelOf(s string) log.Lvl {
	switch s {
	case "debug":
		return log.DEBUG
	case "info":
		return log.INFO
	case "warn":
		return log.WARN
	case "error":
		return log.ERROR
	default:
		return log.OFF
	}
}
