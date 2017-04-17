package k8s

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/golang/glog"

	"k8s.io/client-go/pkg/api"
	certutil "k8s.io/client-go/pkg/util/cert"
	"k8s.io/client-go/rest"
)

// NewClusterConfig create a cluster config from a token
func NewClusterConfig() (*rest.Config, error) {
	svcAccountDirectory := "/var/run/secrets/kubernetes.io/serviceaccount"
	if svcAccountDirectoryVar := os.Getenv("SERVICE_ACCOUNT_DIRECTORY"); svcAccountDirectoryVar != "" {
		svcAccountDirectory = svcAccountDirectoryVar
	}
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	if len(host) == 0 || len(port) == 0 {
		return nil, fmt.Errorf("unable to load in-cluster configuration, KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT must be defined")
	}

	token, err := ioutil.ReadFile(svcAccountDirectory + "/" + api.ServiceAccountTokenKey)
	if err != nil {
		return nil, err
	}
	tlsClientConfig := rest.TLSClientConfig{}
	rootCAFile := svcAccountDirectory + "/" + api.ServiceAccountRootCAKey
	if _, err := certutil.NewPool(rootCAFile); err != nil {
		glog.Errorf("Expected to load root CA config from %s, but got err: %v", rootCAFile, err)
	} else {
		tlsClientConfig.CAFile = rootCAFile
	}

	return &rest.Config{
		// TODO: switch to using cluster DNS.
		Host:            "https://" + net.JoinHostPort(host, port),
		BearerToken:     string(token),
		TLSClientConfig: tlsClientConfig,
	}, nil
}
