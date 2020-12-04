module github.com/redhat-marketplace/redhat-marketplace-operator/v2/metering

go 1.15

require (
	emperror.dev/errors v0.8.0
	github.com/cespare/xxhash v1.1.0
	github.com/go-logr/logr v0.3.0
	github.com/google/wire v0.4.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/openshift/origin v4.1.0+incompatible
	github.com/operator-framework/api v0.3.25
	github.com/petermattis/goid v0.0.0-20180202154549-b0b1615b78e5 // indirect
	github.com/prometheus-operator/prometheus-operator v0.44.0
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.44.0
	github.com/prometheus/client_golang v1.8.0
	github.com/sasha-s/go-deadlock v0.2.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	golang.org/x/net v0.0.0-20201202161906-c7110b5ffcbb
	k8s.io/api v0.19.4
	k8s.io/apimachinery v0.19.4
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
	k8s.io/kube-state-metrics v1.9.7
	sigs.k8s.io/controller-runtime v0.6.4
)
