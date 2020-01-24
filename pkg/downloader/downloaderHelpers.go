package downloader

// NewHelmDownloader provides a downloader for the 'helm' command line utility
func NewHelmDownloader(version string, downloadPath string) Downloader {
	return NewBaseDownloader("https://get.helm.sh/helm-v{VERSION}-{OS}-{ARCH}.tar.gz", GetBaseVersion(version), "helm", downloadPath)
}

// NewKubeDownloader provides a downloader for the 'kubectl' command line utility
func NewKubeDownloader(version string, downloadPath string) Downloader {
	return NewBaseDownloader("https://dl.k8s.io/v{VERSION}/kubernetes-client-{OS}-{ARCH}.tar.gz", GetBaseVersion(version), "kubectl", downloadPath)
}

// NewVaultDownloader provides a downloader for the 'vault' command line utility
func NewVaultDownloader(version string, downloadPath string) Downloader {
	return NewBaseDownloader("https://releases.hashicorp.com/vault/{VERSION}/vault_{VERSION}_{OS}_{ARCH}.zip", GetBaseVersion(version), "vault", downloadPath)
}
