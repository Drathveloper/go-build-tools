package main

import (
	"go-build-tools/scripts/internal/helpers"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	baseURL     = "https://github.com/protocolbuffers/protobuf/releases/download"
	binName     = "protoc"
	toolsDir    = "tools"
	allowedArgs = 3
)

func main() {
	if len(os.Args) != allowedArgs {
		panic("usage: setup-grpc <path> <version>")
	}
	basePath := os.Args[1]
	grpcVersion := os.Args[2]
	toolsAbsDir := filepath.Join(basePath, toolsDir)
	if _, err := os.Stat(filepath.Join(toolsAbsDir, binName)); err == nil {
		os.Exit(0)
	}
	if err := helpers.CreateToolsDir(toolsAbsDir); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	zipName, err := helpers.GetZipName(grpcVersion)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if err = helpers.DownloadZip(baseURL, toolsAbsDir, zipName, grpcVersion); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	contentZipName := strings.TrimSuffix(zipName, ".zip")
	zipPath := filepath.Join(toolsAbsDir, zipName)
	contentPath := filepath.Join(toolsAbsDir, contentZipName)
	if err = helpers.Unzip(zipPath, contentPath); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if err = os.Remove(zipPath); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	execPath := filepath.Join(contentPath, "bin", binName)
	outPath := filepath.Join(toolsAbsDir, binName)
	if err = os.Rename(execPath, outPath); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if err = os.RemoveAll(contentPath); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
