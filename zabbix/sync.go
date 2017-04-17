package zabbix

import (
	"fmt"
	"os"
	"strings"
	"time"

	glob "github.com/ryanuber/go-glob"

	kubernetes "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

type ZabbixSynchronizer struct {
	zabbix     *Zabbix
	k8s        *kubernetes.Clientset
	namespaces []string
}

func NewZabbixSynchronizer(k8s *kubernetes.Clientset) (*ZabbixSynchronizer, error) {

	namespaces := "*-validation,*-production,logging,default,rados-gw"
	if namespacesVar := os.Getenv("NAMESPACES"); namespacesVar != "" {
		namespaces = namespacesVar
	}

	url, user, password, rateURL :=
		os.Getenv("ZABBIX_URL"), os.Getenv("ZABBIX_USER"), os.Getenv("ZABBIX_PASSWORD"), os.Getenv("RATE_SERVICE_URL")
	if len(url) == 0 || len(user) == 0 || len(password) == 0 || len(namespaces) == 0 || len(rateURL) == 0 {
		return nil, fmt.Errorf("unable to load Zabbix configuration, ZABBIX_URL, ZABBIX_USER, ZABBIX_PASSWORD and RATE_SERVICE_URL must be defined")
	}
	return &ZabbixSynchronizer{
		zabbix:     NewZabbix(url, user, password, rateURL),
		k8s:        k8s,
		namespaces: strings.Split(namespaces, ","),
	}, nil
}

// SyncZabbix synchronizes Openshift namespaces and Zabbix items in background every 5 minutes
func (zs *ZabbixSynchronizer) Sync() {
	nsList, err := zs.k8s.Core().Namespaces().List(v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, candidate := range nsList.Items {
		if match(zs.namespaces, candidate.Name) {
			fmt.Printf("Add or update monitoring for namespace %s\n", candidate.Name)
			zs.zabbix.NewOrUpdateMonitoring(candidate.Name)
		}
	}
}

func match(patterns []string, candidate string) bool {
	for _, pattern := range patterns {
		if glob.Glob(pattern, candidate) {
			return true
		}
	}
	return false
}

// SyncZabbix synchronizes Openshift namespaces and Zabbix items in background every 5 minutes
func (zs *ZabbixSynchronizer) Start() {
	ticker := time.NewTicker(5 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// do stuff
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
