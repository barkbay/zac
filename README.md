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
```./zac sync rates```