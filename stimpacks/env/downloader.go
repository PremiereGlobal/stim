package env // import "github.com/PremiereGlobal/stim/stimpacks/env"

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/krolaw/zipstream"

	"github.com/PremiereGlobal/stim/pkg/stimlog"
)

type Downloader interface {
	SetVersion(version string)
	GetVersion() string
	GetDownloadURL() string
	GetBinDir() string
	GetBinPath() string
	GetBinName() string
	GetLogger() stimlog.StimLogger
}

func DownloadPackage(dlr Downloader) error {
	binPath := dlr.GetBinPath()
	urlDL := dlr.GetDownloadURL()
	dlr.GetLogger().Debug("Downloading resource from:{}" + urlDL)
	if data, err := os.Stat(binPath); !os.IsNotExist(err) && data.Size() > 40000 {
		dlr.GetLogger().Debug("Already have file:{}", binPath)
		return nil
	}

	resp, err := http.Get(urlDL)
	if err != nil {
		dlr.GetLogger().Warn(err)
		return err
	}
	defer resp.Body.Close()
	if strings.HasSuffix(urlDL, ".tar.gz") {
		archive, err := gzip.NewReader(resp.Body)
		if err != nil {
			dlr.GetLogger().Warn(err)
			return err
		}
		tr := tar.NewReader(archive)
		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				return err
			}
			if err != nil {
				log.Fatal(err)
			}
			if hdr.FileInfo().Mode().IsRegular() && strings.HasSuffix(hdr.Name, dlr.GetBinName()) {
				out, err := os.Create(binPath)
				if err != nil {
					dlr.GetLogger().Warn(err)
					return err
				}
				defer out.Close()
				_, err = io.Copy(out, tr)
				if err != nil {
					return nil
				}
				os.Chmod(binPath, 0755)
				break
			}
		}
	} else if strings.HasSuffix(urlDL, ".zip") {
		archive := zipstream.NewReader(resp.Body)
		for {
			hdr, err := archive.Next()
			if err == io.EOF {
				return err
			}
			if err != nil {
				log.Fatal(err)
			}
			if hdr.FileInfo().Mode().IsRegular() && strings.HasSuffix(hdr.Name, dlr.GetBinName()) {
				out, err := os.Create(binPath)
				if err != nil {
					dlr.GetLogger().Warn(err)
					return err
				}
				defer out.Close()
				_, err = io.Copy(out, archive)
				if err != nil {
					return nil
				}
				os.Chmod(binPath, 0755)
				break
			}
		}
	}
	return nil
}

func MakeEnvLink(dlr Downloader, envPath string, envName string) error {
	binPath := dlr.GetBinPath()
	finalLinkPath := filepath.FromSlash(envPath + "/" + dlr.GetBinName())
	if data, err := os.Lstat(finalLinkPath); !os.IsNotExist(err) {
		dlr.GetLogger().Debug("{} already exists", binPath)
		if data.Mode()&os.ModeSymlink != 0 {
			dlr.GetLogger().Debug("{} is already a symlink", binPath)
			linkPath, err2 := os.Readlink(finalLinkPath)
			if err2 != nil {
				return err2
			}
			if binPath == linkPath {
				dlr.GetLogger().Debug("{} already pointing to correct location: {}", binPath, linkPath)
				return nil
			} else {
				dlr.GetLogger().Debug("Removing old symlink to:{}", linkPath)
				os.Remove(finalLinkPath)
			}
		} else {
			os.Remove(finalLinkPath)
		}
	}
	os.Symlink(binPath, finalLinkPath)
	return nil
}
