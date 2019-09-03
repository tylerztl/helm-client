package helm

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"
	"helm-client/commons"

	"github.com/prometheus/common/log"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/helm/cmd/helm/installer"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/helm/portforwarder"
	"k8s.io/helm/pkg/kube"
	"k8s.io/helm/pkg/tlsutil"
)

const (
	defaultStableRepositoryUrl = "https://kubernetes-charts.storage.googleapis.com"
	defaultLocalRepositoryURL  = "http://127.0.0.1:8879/charts"
)

var helmClient helm.Interface

func init() {
	log.Info("Initialize Helm client")
	stableRepositoryURL := os.Getenv(string(commons.StableRepositoryURL))
	if stableRepositoryURL == "" {
		stableRepositoryURL = defaultStableRepositoryUrl
	}
	localRepositoryURL := os.Getenv(string(commons.LocalRepositoryURL))
	if localRepositoryURL == "" {
		localRepositoryURL = defaultLocalRepositoryURL
	}
	skipRefresh := true
	skipRefreshStr := os.Getenv(string(commons.SkipRefresh))
	if skipRefreshStr != "" {
		var err error
		skipRefresh, err = strconv.ParseBool(skipRefreshStr)
		if err != nil {
			panic("skipRefresh ParseBool error")
		}
	}

	if err := installer.Initialize(commons.Settings.Home, bytes.NewBuffer(nil), skipRefresh, commons.Settings, stableRepositoryURL, localRepositoryURL); err != nil {
		panic(fmt.Errorf("error initializing: %s", err))
	}

	// establish a connection to Tiller now that we've effectively guaranteed it's available
	if err := setupConnection(); err != nil {
		panic("Failed to connect tiller server")
	}

	chConnect := make(chan bool)
	go func() {
		helmClient = newClient()
		if err := helmClient.PingTiller(); err != nil {
			panic(fmt.Sprintf("Could not ping tiller server: %s", err))
		} else {
			chConnect <- true
		}
	}()

	tick := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-tick.C:
			panic("Connect to tiller server timeout")
		case <-chConnect:
			log.Info("Connect to tiller server successful")
			return
		}
	}
}

func GetClient() helm.Interface {
	return helmClient
}

func setupConnection() error {
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
		log.Info(fmt.Sprintf("Created tunnel using local port: '%d'", tillerTunnel.Local))
	}

	// Set up the gRPC config.
	log.Info(fmt.Sprintf("SERVER: %q", commons.Settings.TillerHost))

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
		log.Debug("Host=%q, Key=%q, Cert=%q, CA=%q\n", settings.TLSServerName, settings.TLSKeyFile, settings.TLSCertFile, settings.TLSCaCertFile)
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
