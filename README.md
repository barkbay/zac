# ZAC - Zabbix Automatic Configuration for Kubernetes

# Warning events rate alerting
## The rate server
```zac server``` computes rate of events of type "Warning" in a Kubernetes cluster

It has been designed to be embedded directly inside a POD.

Do not forget to add view role to the service account :

```
oc adm policy add-cluster-role-to-user view system:serviceaccount:mynamespace:default
```

Where ```mynamespace``` is the namespace where the r8-alter pod is living.

## Create Web Scenarios in Zabbix
```./zac sync rates``` creates Web Scenarios and triggers in Zabbix

Example :

```
> $ NAMESPACES="default,*-dev,jenkins*" RATE_SERVICE_URL=http://192.168.99.1:8080/ ZABBIX_URL=http://localhost:8070 ZABBIX_USER=admin ZABBIX_PASSWORD=zabbix SERVICE_ACCOUNT_DIRECTORY=. KUBERNETES_SERVICE_HOST=192.168.99.100 KUBERNETES_SERVICE_PORT=8443 ./zac sync rates
==> Add or Update Zabbix web scenario and triggers
Add or update monitoring for namespace default
Add or update monitoring for namespace myproject-dev
Add or update monitoring for namespace jenkins-foo
Add or update monitoring for namespace jenkins2
```

# Environment variables

You must set these variables :

* **NAMESPACES**
  * List of namespaces to monitor, separated with commas
* **RATE_SERVICE_URL**
  * Where the rate service is hosted. It is used to configurure the URL in the step of the Web Scenario
* **ZABBIX_URL**
  * Zabbix URL
* **ZABBIX_USER**
  * Zabbix username to use for the API calls
* **ZABBIX_PASSWORD**
  * Zabbix username's password
* **SERVICE_ACCOUNT_DIRECTORY**
  * Where to find the Kubernetes files ```ca.crt``` and ```token``` (not required if inside a POD)
* **KUBERNETES_SERVICE_HOST**
  * Kubernetes host
* **KUBERNETES_SERVICE_PORT**
  * Kubernetes port