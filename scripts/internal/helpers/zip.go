package helpers

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const (
	linuxOS   = "linux"
	windowsOS = "windows"
	macOS     = "darwin"

	filePermission = 0o750
)

// ErrUnsupportedOS is returned when the OS is not supported.
var ErrUnsupportedOS = errors.New("unsupported OS")

// ErrUnsupportedArch is returned when the architecture is not supported.
var ErrUnsupportedArch = errors.New("unsupported architecture")

// CreateToolsDir creates the tools' directory.
func CreateToolsDir(dir string) error {
	if err := os.MkdirAll(dir, filePermission); err != nil {
		return fmt.Errorf("error creating tools directory: %w", err)
	}
	return nil
}

// GetZipName returns the name of the zip file for the given grpc version.
func GetZipName(grpcVersion string) (string, error) {
	switch runtime.GOOS {
	case linuxOS:
		if runtime.GOARCH == "amd64" {
			return fmt.Sprintf("protoc-%s-linux-x86_64.zip", grpcVersion), nil
		}
		if runtime.GOARCH == "arm64" {
			return fmt.Sprintf("protoc-%s-linux-aarch_64.zip", grpcVersion), nil
		}
		return "", fmt.Errorf("get zip name failed for %s/%s failed: %w", runtime.GOOS, runtime.GOARCH, ErrUnsupportedArch)
	case macOS:
		if runtime.GOARCH == "amd64" {
			return fmt.Sprintf("protoc-%s-osx-x86_64.zip", grpcVersion), nil
		}
		if runtime.GOARCH == "arm64" {
			return fmt.Sprintf("protoc-%s-osx-aarch_64.zip", grpcVersion), nil
		}
		return "", fmt.Errorf("get zip name failed for %s/%s failed: %w", runtime.GOOS, runtime.GOARCH, ErrUnsupportedArch)
	case windowsOS:
		return fmt.Sprintf("protoc-%s-win64.zip", grpcVersion), nil
	default:
		return "", fmt.Errorf("get zip name failed for %s/%s failed: %w", runtime.GOOS, runtime.GOARCH, ErrUnsupportedOS)
	}
}

// DownloadZip downloads a zip file from the baseURL.
func DownloadZip(baseURL, baseDir, zipName, grpcVersion string) error {
	url := fmt.Sprintf("%s/v%s/%s", baseURL, grpcVersion, zipName)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error downloading %s: %w", url, err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error downloading %s: %w", url, err)
	}
	defer resp.Body.Close() //nolint: errcheck
	zipPath := filepath.Join(baseDir, zipName)
	out, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("error downloading %s: %w", url, err)
	}
	defer out.Close() //nolint: errcheck
	if _, err = io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("error downloading %s: %w", url, err)
	}
	return nil
}

// Unzip unzips the zip file to the contentPath.
func Unzip(zipPath, contentPath string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("error unzipping %s: %w", zipPath, err)
	}
	defer reader.Close() //nolint: errcheck
	for _, f := range reader.File {
		if err = unzipFile(f, contentPath); err != nil {
			return fmt.Errorf("error unzipping %s: %w", zipPath, err)
		}
	}
	return nil
}

func unzipFile(file *zip.File, dst string) error {
	fpath := filepath.Join(dst, file.Name)
	if file.FileInfo().IsDir() {
		fmt.Println("creating dir", fpath)
		if err := os.MkdirAll(fpath, file.Mode()); err != nil {
			return fmt.Errorf("error unzipping individual file %s: %w", filepath.Dir(fpath), err)
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(fpath), filePermission); err != nil {
		return fmt.Errorf("error unzipping individual file %s: %w", filepath.Dir(fpath), err)
	}
	inFile, err := file.Open()
	if err != nil {
		return fmt.Errorf("error unzipping individual file %s: %w", filepath.Dir(fpath), err)
	}
	defer inFile.Close()
	outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return fmt.Errorf("error unzipping individual file %s: %w", filepath.Dir(fpath), err)
	}
	defer outFile.Close()
	if _, err = io.Copy(outFile, inFile); err != nil {
		return fmt.Errorf("error unzipping individual file %s: %w", filepath.Dir(fpath), err)
	}
	return nil
}
