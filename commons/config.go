package commons

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	helm_env "k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/helmpath"
)

var Settings helm_env.EnvSettings

type AppEnv string

const (
	TillerHost          = AppEnv("TILLER_HOST")
	HelmHome            = AppEnv("HELM_HOME")
	SkipRefresh         = AppEnv("SKIP_REFRESH") // do not refresh (download) the local repository cache
	StableRepositoryURL = AppEnv("STABLE_REPOSITORY_URL")
	LocalRepositoryURL  = AppEnv("LOCAL_REPOSITORY_URL")
)

func init() {
	flags := pflag.NewFlagSet("helm-client", pflag.ContinueOnError)
	Settings.AddFlags(flags)
	Settings.AddFlagsTLS(flags)
	Settings.Init(flags)
	Settings.InitTLS(flags)
	initAppEnv()
}

func initAppEnv() {
	tillerHost := os.Getenv(string(TillerHost))
	if tillerHost != "" {
		Settings.TillerHost = tillerHost
	}
	helmHome := os.Getenv(string(HelmHome))
	if helmHome != "" {
		if fi, err := os.Stat(helmHome); err != nil {
			panic("Invalid helm home")
		} else if !fi.IsDir() {
			panic(fmt.Sprintf("helm home [%s] must be a directory", helmHome))
		}
		Settings.Home = helmpath.Home(helmHome)
	}
}

func GetConfig() helm_env.EnvSettings {
	return Settings
}
