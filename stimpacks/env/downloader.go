package env // import "github.com/PremiereGlobal/stim/stimpacks/env"

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
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
	GetBinBaseName() string
	GetLogger() stimlog.StimLogger
}

type baseDownloader struct {
	version, name, path string
	url                 StringReplacer
}

func NewBaseDownloader(url, version, name, path string) Downloader {
	return &baseDownloader{
		url:     StringReplacer(url),
		version: version,
		name:    name,
		path:    path,
	}
}

func (bd *baseDownloader) SetVersion(version string) {
	bd.version = version
}
func (bd *baseDownloader) GetVersion() string {
	return bd.version
}
func (bd *baseDownloader) GetDownloadURL() string {
	return bd.url.ReplaceAll("{VERSION}", bd.version).
		ReplaceAll("{OS}", runtime.GOOS).
		ReplaceAll("{ARCH}", runtime.GOARCH).
		ReplaceAll("{NAME}", bd.name).
		String()

	// Go sucks
	//	return strings.ReplaceAll(
	//		strings.ReplaceAll(
	//			strings.ReplaceAll(
	//				strings.ReplaceAll(bd.url, "{VERSION}", bd.version),
	//				"{NAME}", bd.name),
	//			"{OS}", runtime.GOOS),
	//		"{ARCH}", runtime.GOARCH)
}
func (bd *baseDownloader) GetBinBaseName() string {
	return bd.name
}
func (bd *baseDownloader) GetBinName() string {
	return bd.name + "-" + bd.version
}
func (bd *baseDownloader) GetLogger() stimlog.StimLogger {
	return stimlog.GetLogger()
}
func (bd *baseDownloader) GetBinDir() string {
	return bd.GetBinPath()
}
func (bd *baseDownloader) GetBinPath() string {
	return filepath.Join(bd.path, bd.GetBinName())
}

func NewKubeDownloader(version string, envBinPath string) Downloader {
	if !strings.HasPrefix(version, "v") {
		nv := "v" + version
		stimlog.GetLogger().Info("Kube version set as:'{}', it MUST start with a 'v', changing to `{}`", version, nv)
		version = nv
	}
	return NewBaseDownloader("https://dl.k8s.io/{VERSION}/kubernetes-client-{OS}-{ARCH}.tar.gz", version, "kubectl", envBinPath)
}

func NewVaultDownloader(version string, envBinPath string) Downloader {
	if strings.HasPrefix(version, "v") {
		nv := version[1:]
		stimlog.GetLogger().Info("Vault version set as:'{}', it DoesNot start with a 'v', changing to `{}`", version, nv)
		version = nv
	}
	return NewBaseDownloader("https://releases.hashicorp.com/vault/{VERSION}/vault_{VERSION}_{OS}_{ARCH}.zip", version, "vault", envBinPath)
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
			if hdr.FileInfo().Mode().IsRegular() && strings.HasSuffix(hdr.Name, dlr.GetBinBaseName()) {
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
			if hdr.FileInfo().Mode().IsRegular() && strings.HasSuffix(hdr.Name, dlr.GetBinBaseName()) {
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

type StringReplacer string

func (sr StringReplacer) ReplaceAll(replace, with string) StringReplacer {
	return StringReplacer(strings.ReplaceAll(string(sr), replace, with))
}
func (sr StringReplacer) String() string {
	return string(sr)
}
