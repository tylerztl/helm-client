package commons

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/prometheus/common/log"
	"k8s.io/helm/pkg/downloader"
	"k8s.io/helm/pkg/getter"
	"k8s.io/helm/pkg/helm/helmpath"
	"k8s.io/helm/pkg/kube"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/repo"
	"k8s.io/helm/pkg/timeconv"
)

// define the helper functions

// locateChartPath looks for a chart directory in known places, and returns either the full path or an error.
//
// This does not ensure that the chart is well-formed; only that the requested filename exists.
//
// Order of resolution:
// - current working directory
// - if path is absolute or begins with '.', error out here
// - chart repos in $HELM_HOME
// - URL
//
// If 'verify' is true, this will attempt to also verify the chart.
func LocateChartPath(repoURL, username, password, name, version string, verify bool, keyring,
certFile, keyFile, caFile string) (string, error) {
	name = strings.TrimSpace(name)
	version = strings.TrimSpace(version)
	if fi, err := os.Stat(name); err == nil {
		abs, err := filepath.Abs(name)
		if err != nil {
			return abs, err
		}
		if verify {
			if fi.IsDir() {
				return "", errors.New("cannot verify a directory")
			}
			if _, err := downloader.VerifyChart(abs, keyring); err != nil {
				return "", err
			}
		}
		return abs, nil
	}
	if filepath.IsAbs(name) || strings.HasPrefix(name, ".") {
		return name, fmt.Errorf("path %q not found", name)
	}

	crepo := filepath.Join(Settings.Home.Repository(), name)
	if _, err := os.Stat(crepo); err == nil {
		return filepath.Abs(crepo)
	}
	dl := downloader.ChartDownloader{
		HelmHome: Settings.Home,
		Out:      os.Stdout,
		Keyring:  keyring,
		Getters:  getter.All(Settings),
		Username: username,
		Password: password,
	}
	if verify {
		dl.Verify = downloader.VerifyAlways
	}
	if repoURL != "" {
		chartURL, err := repo.FindChartInAuthRepoURL(repoURL, username, password, name, version,
			certFile, keyFile, caFile, getter.All(Settings))
		if err != nil {
			return "", err
		}
		name = chartURL
	}

	if _, err := os.Stat(Settings.Home.Archive()); os.IsNotExist(err) {
		os.MkdirAll(Settings.Home.Archive(), 0744)
	}

	filename, _, err := dl.DownloadTo(name, version, Settings.Home.Archive())
	if err == nil {
		lname, err := filepath.Abs(filename)
		if err != nil {
			return filename, err
		}
		log.Debug(fmt.Sprintf("Fetched %s to %s\n", name, filename))
		return lname, nil
	}

	return filename, fmt.Errorf("failed to download %q due to %v (hint: running `helm repo update` may help)", name, err)
}

func AddRepository(name, url, username, password string, home helmpath.Home, certFile, keyFile, caFile string, noUpdate bool) error {
	f, err := repo.LoadRepositoriesFile(home.RepositoryFile())
	if err != nil {
		return err
	}

	if noUpdate && f.Has(name) {
		return fmt.Errorf("repository name (%s) already exists, please specify a different name", name)
	}

	cif := home.CacheIndex(name)
	c := repo.Entry{
		Name:     name,
		Cache:    cif,
		URL:      url,
		Username: username,
		Password: password,
		CertFile: certFile,
		KeyFile:  keyFile,
		CAFile:   caFile,
	}

	r, err := repo.NewChartRepository(&c, getter.All(GetConfig()))
	if err != nil {
		return err
	}

	if err := r.DownloadIndexFile(home.Cache()); err != nil {
		return fmt.Errorf("Looks like %q is not a valid chart repository or cannot be reached: %s", url, err.Error())
	}

	f.Update(&c)

	return f.WriteFile(home.RepositoryFile(), 0644)
}

func RemoveRepoLine(name string, home helmpath.Home) error {
	repoFile := home.RepositoryFile()
	r, err := repo.LoadRepositoriesFile(repoFile)
	if err != nil {
		return err
	}

	if !r.Remove(name) {
		return fmt.Errorf("no repo named %q found", name)
	}
	if err := r.WriteFile(repoFile, 0644); err != nil {
		return err
	}

	if err := removeRepoCache(name, home); err != nil {
		return err
	}

	log.Debug(fmt.Sprintf("%q has been removed from your repositories", name))

	return nil
}

func removeRepoCache(name string, home helmpath.Home) error {
	if _, err := os.Stat(home.CacheIndex(name)); err == nil {
		err = os.Remove(home.CacheIndex(name))
		if err != nil {
			return err
		}
	}
	return nil
}

// filterList returns a list scrubbed of old releases.
func FilterList(rels []*release.Release) []*release.Release {
	idx := map[string]int32{}

	for _, r := range rels {
		name, version := r.GetName(), r.GetVersion()
		if max, ok := idx[name]; ok {
			// check if we have a greater version already
			if max > version {
				continue
			}
		}
		idx[name] = version
	}

	uniq := make([]*release.Release, 0, len(idx))
	for _, r := range rels {
		if idx[r.GetName()] == r.GetVersion() {
			uniq = append(uniq, r)
		}
	}
	return uniq
}

func GetListResult(rels []*release.Release, next string) *ListResult {
	listReleases := []ListRelease{}
	for _, r := range rels {
		md := r.GetChart().GetMetadata()
		t := "-"
		if tspb := r.GetInfo().GetLastDeployed(); tspb != nil {
			t = timeconv.String(tspb)
		}

		lr := ListRelease{
			Name:       r.GetName(),
			Revision:   r.GetVersion(),
			Updated:    t,
			Status:     r.GetInfo().GetStatus().GetCode().String(),
			Chart:      fmt.Sprintf("%s-%s", md.GetName(), md.GetVersion()),
			AppVersion: md.GetAppVersion(),
			Namespace:  r.GetNamespace(),
		}
		listReleases = append(listReleases, lr)
	}

	return &ListResult{
		Releases: listReleases,
		Next:     next,
	}
}

func DefaultNamespace() string {
	if ns, _, err := kube.GetConfig(Settings.KubeContext, Settings.KubeConfig).Namespace(); err == nil {
		return ns
	}
	return "default"
}

func MakeReleaseResource(release *release.Release) *ReleaseResource {
	if release == nil {
		return nil
	}
	return &ReleaseResource{
		ChartName:    release.Chart.Metadata.Name,
		ChartVersion: release.Chart.Metadata.Version,
		ChartIcon:    release.Chart.Metadata.Icon,
		Updated:      timeconv.String(release.Info.LastDeployed),
		Name:         release.Name,
		Namespace:    release.Namespace,
		Status:       release.Info.Status.Code.String(),
	}
}

func MakeReleaseExtendedResource(release *release.Release) *ReleaseExtended {
	if release == nil {
		return nil
	}
	return &ReleaseExtended{
		ChartName:    release.Chart.Metadata.Name,
		ChartVersion: release.Chart.Metadata.Version,
		ChartIcon:    release.Chart.Metadata.Icon,
		Updated:      timeconv.String(release.Info.LastDeployed),
		Name:         release.Name,
		Namespace:    release.Namespace,
		Status:       release.Info.Status.Code.String(),
		Resources:    release.Info.Status.Resources,
		Notes:        release.Info.Status.Notes,
	}
}
