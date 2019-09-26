package env // import "github.com/PremiereGlobal/stim/stimpacks/env"

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/PremiereGlobal/stim/pkg/stimlog"
)

type KubeDownloader struct {
	version string
	env     Env
}

const kubeDownloadURL string = "https://dl.k8s.io"
const kubeFileName string = "kubectl"

func (kd *KubeDownloader) SetVersion(version string) {
	kd.version = version
}

func (kd *KubeDownloader) GetVersion() string {
	return kd.version
}

func (kd *KubeDownloader) GetLogger() stimlog.StimLogger {
	return kd.env.stim.GetLogger()
}

func (kd *KubeDownloader) GetDownloadURL() string {
	return fmt.Sprintf("%s/%s/kubernetes-client-%s-%s.tar.gz", kubeDownloadURL, kd.version, runtime.GOOS, runtime.GOARCH)
}

func (kd *KubeDownloader) GetBinName() string {
	return kubeFileName
}

func (kd *KubeDownloader) GetBinPath() string {
	return filepath.FromSlash(fmt.Sprintf("%s/kubectl-%s", kd.GetBinDir(), kd.version))
}

func (kd *KubeDownloader) GetBinDir() string {
	return filepath.FromSlash(kd.env.GetEnvBinDir())
}

func (kd *KubeDownloader) makeKubernetesDir() error {
	if _, err := os.Stat(kd.GetBinDir()); os.IsNotExist(err) {
		err := os.MkdirAll(kd.GetBinDir(), 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
