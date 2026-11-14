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
	"strings"
)

const (
	linuxOS   = "linux"
	windowsOS = "windows"
	macOS     = "darwin"

	filePermission = 0o750

	// maxUnzippedFileSize limits the size of each extracted file to protect against decompression bombs.
	maxUnzippedFileSize = 512 << 20 // 512 MiB
)

// ErrUnsupportedOS is returned when the OS is not supported.
var ErrUnsupportedOS = errors.New("unsupported OS")

// ErrUnsupportedArch is returned when the architecture is not supported.
var ErrUnsupportedArch = errors.New("unsupported architecture")

// ErrUnexpectedStatus is returned when the download server responds with a non-200 status.
var ErrUnexpectedStatus = errors.New("unexpected HTTP status")

// ErrIllegalFilePath is returned when a zip entry would be extracted outside the destination directory.
var ErrIllegalFilePath = errors.New("illegal file path in zip archive")

// ErrFileTooLarge is returned when a zip entry exceeds maxUnzippedFileSize once extracted.
var ErrFileTooLarge = errors.New("unzipped file too large")

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
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error downloading %s: %w: %s", url, ErrUnexpectedStatus, resp.Status)
	}
	zipPath := filepath.Join(baseDir, zipName)
	out, err := os.Create(filepath.Clean(zipPath))
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
	fpath, err := safeJoin(dst, file.Name)
	if err != nil {
		return fmt.Errorf("error unzipping individual file %s: %w", file.Name, err)
	}
	if file.FileInfo().IsDir() {
		if err = os.MkdirAll(fpath, file.Mode()); err != nil {
			return fmt.Errorf("error unzipping individual file %s: %w", fpath, err)
		}
		return nil
	}
	if err = os.MkdirAll(filepath.Dir(fpath), filePermission); err != nil {
		return fmt.Errorf("error unzipping individual file %s: %w", fpath, err)
	}
	if err = extractFile(file, fpath); err != nil {
		return fmt.Errorf("error unzipping individual file %s: %w", fpath, err)
	}
	return nil
}

// safeJoin joins dst and name, rejecting names that would escape dst (zip slip).
func safeJoin(dst, name string) (string, error) {
	if !filepath.IsLocal(name) {
		return "", ErrIllegalFilePath
	}
	fpath := filepath.Join(dst, name)
	if !strings.HasPrefix(fpath, filepath.Clean(dst)+string(os.PathSeparator)) {
		return "", ErrIllegalFilePath
	}
	return fpath, nil
}

// extractFile copies the content of a zip entry to fpath, limiting its size to maxUnzippedFileSize.
func extractFile(file *zip.File, fpath string) (err error) {
	inFile, err := file.Open()
	if err != nil {
		return fmt.Errorf("error opening zip entry: %w", err)
	}
	defer inFile.Close() //nolint: errcheck
	outFile, err := os.OpenFile(filepath.Clean(fpath), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer func() {
		if closeErr := outFile.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("error closing file: %w", closeErr)
		}
	}()
	written, err := io.CopyN(outFile, inFile, maxUnzippedFileSize+1)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("error writing file: %w", err)
	}
	if written > maxUnzippedFileSize {
		return ErrFileTooLarge
	}
	return nil
}
