package env // import "github.com/PremiereGlobal/stim/stimpacks/env"

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/PremiereGlobal/stim/pkg/stimlog"
)

type VaultDownloader struct {
	version string
	env     Env
}

const vaultDownloadURL string = "https://releases.hashicorp.com/vault"
const vaultFileName string = "vault"

func (vd *VaultDownloader) SetVersion(version string) {
	vd.version = version
}

func (vd *VaultDownloader) GetVersion() string {
	return vd.version
}

func (vd *VaultDownloader) GetLogger() stimlog.StimLogger {
	return vd.env.stim.GetLogger()
}

func (vd *VaultDownloader) GetDownloadURL() string {
	return fmt.Sprintf("%s/%s/vault_%s_%s_%s.zip", vaultDownloadURL, vd.version, vd.version, runtime.GOOS, runtime.GOARCH)
}

func (vd *VaultDownloader) GetBinName() string {
	return vaultFileName
}

func (vd *VaultDownloader) GetBinPath() string {
	return filepath.FromSlash(fmt.Sprintf("%s/vault-%s", vd.GetBinDir(), vd.version))
}

func (vd *VaultDownloader) GetBinDir() string {
	return filepath.FromSlash(vd.env.GetEnvBinDir())
}

func (vd *VaultDownloader) makeKubernetesDir() error {
	if _, err := os.Stat(vd.GetBinDir()); os.IsNotExist(err) {
		err := os.MkdirAll(vd.GetBinDir(), 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
