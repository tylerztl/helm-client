package helm

import (
	"fmt"
	"os"

	// Import to initialize client auth plugins.
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/helm/portforwarder"
	"k8s.io/helm/pkg/kube"
	"k8s.io/helm/pkg/tlsutil"
	"zig-helm/commons"
)

const (
	envTillerHost = "TILLER_HOST"
)

var helmClient helm.Interface

func init() {
	// establish a connection to Tiller now that we've effectively guaranteed it's available
	if err := setupConnection(); err != nil {
		panic("Failed to connect tiller server")
	}
	helmClient = newClient()
	//if err := helmClient.PingTiller(); err != nil {
	//	panic(fmt.Sprintf("could not ping Tiller: %s", err))
	//}
}

func GetClient() helm.Interface {
	return helmClient
}

func setupConnection() error {
	tillerHost := os.Getenv(envTillerHost)
	commons.Settings.TillerHost = tillerHost

	if commons.Settings.TillerHost == "" {
		config, client, err := getKubeClient(commons.Settings.KubeContext, commons.Settings.KubeConfig)
		if err != nil {
			return err
		}

		tillerTunnel, err := portforwarder.New(commons.Settings.TillerNamespace, client, config)
		if err != nil {
			return err
		}

		commons.Settings.TillerHost = fmt.Sprintf("127.0.0.1:%d", tillerTunnel.Local)
		fmt.Printf("Created tunnel using local port: '%d'\n", tillerTunnel.Local)
	}

	// Set up the gRPC config.
	fmt.Printf("SERVER: %q\n", commons.Settings.TillerHost)

	// Plugin support.
	return nil
}

// configForContext creates a Kubernetes REST client configuration for a given kubeconfig context.
func configForContext(context string, kubeconfig string) (*rest.Config, error) {
	config, err := kube.GetConfig(context, kubeconfig).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get Kubernetes config for context %q: %s", context, err)
	}
	return config, nil
}

// getKubeClient creates a Kubernetes config and client for a given kubeconfig context.
func getKubeClient(context string, kubeconfig string) (*rest.Config, kubernetes.Interface, error) {
	config, err := configForContext(context, kubeconfig)
	if err != nil {
		return nil, nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get Kubernetes client: %s", err)
	}
	return config, client, nil
}

func newClient() helm.Interface {
	settings := commons.GetConfig()
	options := []helm.Option{helm.Host(settings.TillerHost), helm.ConnectTimeout(settings.TillerConnectionTimeout)}

	if settings.TLSVerify || settings.TLSEnable {
		fmt.Printf("Host=%q, Key=%q, Cert=%q, CA=%q\n", settings.TLSServerName, settings.TLSKeyFile, settings.TLSCertFile, settings.TLSCaCertFile)
		tlsopts := tlsutil.Options{
			ServerName:         settings.TLSServerName,
			KeyFile:            settings.TLSKeyFile,
			CertFile:           settings.TLSCertFile,
			InsecureSkipVerify: true,
		}
		if settings.TLSVerify {
			tlsopts.CaCertFile = settings.TLSCaCertFile
			tlsopts.InsecureSkipVerify = false
		}
		tlscfg, err := tlsutil.ClientConfig(tlsopts)
		if err != nil {
			panic(err)
		}
		options = append(options, helm.WithTLS(tlscfg))
	}
	return helm.NewClient(options...)
}
