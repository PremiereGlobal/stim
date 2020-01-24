package downloader

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/PremiereGlobal/stim/pkg/utils"
	"github.com/krolaw/zipstream"
)

// Downloader represents the Downloader type (duh)
type Downloader interface {
	Download() (DownloadResult, error)
	SetVersion(version string)
	GetVersion() string
	GetDownloadURL() string
	GetBinPath() string
	GetBinName() string
	GetBinBaseName() string
}

type baseDownloader struct {
	version, name, path string
	url                 utils.StringReplacer
}

type DownloadResult struct {
	RenderedURL      string
	FileExists       bool
	DownloadDuration time.Duration
}

// NewBaseDownloader returns a New baseDownloader
// url = URL template for the download
// version = version to download
// name = name of the binary file
// path = where the binary file should end up
func NewBaseDownloader(url, version, name, path string) Downloader {
	d := &baseDownloader{
		url:  utils.StringReplacer(url),
		name: name,
		path: path,
	}

	d.SetVersion(version)

	return d
}

// SetVersion sets the version to download
func (bd *baseDownloader) SetVersion(version string) {
	bd.version = GetBaseVersion(version)
}

// SetVersion gets the version to be downloaded
func (bd *baseDownloader) GetVersion() string {
	return bd.version
}

// GetDownloadURL returns the constructed download url
func (bd *baseDownloader) GetDownloadURL() string {
	return bd.url.ReplaceAll("{VERSION}", bd.version).
		ReplaceAll("{OS}", runtime.GOOS).
		ReplaceAll("{ARCH}", runtime.GOARCH).
		ReplaceAll("{NAME}", bd.name).
		String()
}

// GetBinBaseName returns the base name of the binary (ex. stim)
func (bd *baseDownloader) GetBinBaseName() string {
	return bd.name
}

// GetBinName returns the name of the binary + version (ex. stim-1.2)
func (bd *baseDownloader) GetBinName() string {
	return bd.name + "-v" + bd.version
}

// GetBinPath returns the full path to the binary
func (bd *baseDownloader) GetBinPath() string {
	return filepath.Join(bd.path, bd.GetBinName())
}

// Download downloads the file and move it to the appropriate path
func (bd *baseDownloader) Download() (DownloadResult, error) {
	binPath := bd.GetBinPath()
	urlDL := bd.GetDownloadURL()
	result := DownloadResult{}
	result.RenderedURL = urlDL

	if data, err := os.Stat(binPath); !os.IsNotExist(err) && data.Size() > 40000 {

		// Ensure the file is executable (by the user)
		if data.Mode().Perm()&0001 == 0 {
			os.Chmod(binPath, 0755)
		}
		result.FileExists = true
		return result, nil
	}

	start := time.Now()
	resp, err := http.Get(urlDL)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	result.DownloadDuration = time.Since(start)

	if strings.HasSuffix(urlDL, ".tar.gz") {
		archive, err := gzip.NewReader(resp.Body)
		if err != nil {
			return result, err
		}
		tr := tar.NewReader(archive)
		for {
			hdr, err := tr.Next()
			if err != nil {
				return result, err
			}

			if hdr.FileInfo().Mode().IsRegular() && strings.HasSuffix(hdr.Name, bd.GetBinBaseName()) {
				out, err := os.Create(binPath)
				if err != nil {
					return result, err
				}
				defer out.Close()
				_, err = io.Copy(out, tr)
				if err != nil {
					return result, err
				}
				os.Chmod(binPath, 0755)
				break
			}
		}
	} else if strings.HasSuffix(urlDL, ".zip") {
		archive := zipstream.NewReader(resp.Body)
		for {
			hdr, err := archive.Next()
			if err != nil {
				return result, err
			}
			if hdr.FileInfo().Mode().IsRegular() && strings.HasSuffix(hdr.Name, bd.GetBinBaseName()) {
				out, err := os.Create(binPath)
				if err != nil {
					return result, err
				}
				defer out.Close()
				_, err = io.Copy(out, archive)
				if err != nil {
					return result, err
				}
				os.Chmod(binPath, 0755)
				break
			}
		}
	}
	return result, nil
}
