package commons

import (
	"github.com/spf13/pflag"
	helm_env "k8s.io/helm/pkg/helm/environment"
)

var Settings helm_env.EnvSettings

func init() {
	flags := pflag.NewFlagSet("zig-helm", pflag.ContinueOnError)
	Settings.AddFlags(flags)
	Settings.AddFlagsTLS(flags)
	Settings.Init(flags)
	Settings.InitTLS(flags)
}

func GetConfig() helm_env.EnvSettings {
	return Settings
}
