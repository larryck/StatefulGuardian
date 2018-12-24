
package app

import (
	"context"
	"math/rand"
	"net/http"
	"sync"
	"time"
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	controller "statefulguardian/pkg/controllers"
	clientset "statefulguardian/pkg/generated/clientset/versioned"
	informers "statefulguardian/pkg/generated/informers/externalversions"
	//operatoropts "statefulguardian/pkg/options/operator"
	//metrics "statefulguardian/pkg/util/metrics"
	//signals "statefulguardian/pkg/util/signals"
)

const (
	metricsEndpoint = "0.0.0.0:8080"
)

const (
	SGOperatorName = "sg-operator"
	SGOperatorNs = "sg-operator"
)


type SgOpts struct {
	// KubeConfig is the path to a kubeconfig file, specifying how to connect to
	// the API server.
	KubeConfig string `yaml:"kubeconfig"`

	// Master is the address of the Kubernetes API server (overrides any value
	// in kubeconfig).
	Master string `yaml:"master"`

	// Namespace is the (optional) namespace in which the MySQL operator will
	// manage MySQL Clusters. Defaults to metav1.NamespaceAll.
	Namespace string `yaml:"namespace"`

	// Hostname of the pod the operator is running in.
	Hostname string `yaml:"hostname"`

	// minResyncPeriod is the resync period in reflectors; will be random
	// between minResyncPeriod and 2*minResyncPeriod.
	MinResyncPeriod metav1.Duration `yaml:"minResyncPeriod"`
}

func NewSgOpts() (*SgOpts, error) {
	hostName, err := os.Hostname()
	if err!= nil {
		return nil, err
	}

	sgo := SgOpts {
		KubeConfig: controller.KUBECONFIG_PATH,
		Namespace: metav1.NamespaceAll,
		Hostname: hostName,
		Master: "",
		MinResyncPeriod: metav1.Duration{Duration: 12 * time.Hour} ,
	}

	return &sgo,nil

}



// resyncPeriod computes the time interval a shared informer waits before
// resyncing with the api server.
func resyncPeriod(s *SgOpts) func() time.Duration {
	return func() time.Duration {
		factor := rand.Float64() + 1
		return time.Duration(float64(s.MinResyncPeriod.Nanoseconds()) * factor)
	}
}

// Run starts the mysql-operator controllers. This should never exit.
func Run(s *SgOpts) error {
	kubeconfig, err := clientcmd.BuildConfigFromFlags(s.Master, s.KubeConfig)
	if err != nil {
		return err
	}

	// Initialise the operator metrics.
	controller.RegisterPodName(s.Hostname)
	controller.RegisterMetrics()
	http.Handle("/metrics", prometheus.Handler())
	go http.ListenAndServe(metricsEndpoint, nil)

	//ctx, cancelFunc := context.WithCancel(context.Background())
	ctx, _ := context.WithCancel(context.Background())

	// Set up signals so we handle the first shutdown signal gracefully.
	//signals.SetupSignalHandler(cancelFunc)

	kubeClient := kubernetes.NewForConfigOrDie(kubeconfig)
	sgClient := clientset.NewForConfigOrDie(kubeconfig)

	// Shared informers (non namespace specific).
	sgInformerFactory := informers.NewFilteredSharedInformerFactory(sgClient, resyncPeriod(s)(), s.Namespace, nil)
	kubeInformerFactory := kubeinformers.NewFilteredSharedInformerFactory(kubeClient, resyncPeriod(s)(), s.Namespace, nil)

	//get the current sg operator deployment
	deploy, err:= kubeClient.AppsV1beta1().Deployments(SGOperatorNs).Get(SGOperatorName, metaV1.GetOptions{})
	if err!=nil {
		fmt.Println("Get sg operator deployment error")
		return err
	}


	var wg sync.WaitGroup

	cont := controller.NewSgController(
		//*s,
		sgClient,
		kubeClient,
		sgInformerFactory.Statefulguardian().V1alpha1().Statefulguardians(),
		kubeInformerFactory.Apps().V1beta1().StatefulSets(),
		kubeInformerFactory.Core().V1().Pods(),
		kubeInformerFactory.Core().V1().Services(),
		30*time.Second,
		s.Namespace,
		ctx,
		deploy,
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		cont.Run(ctx, 5)
	}()

	// Shared informers have to be started after ALL controllers.
	go kubeInformerFactory.Start(ctx.Done())
	go sgInformerFactory.Start(ctx.Done())

	<-ctx.Done()

	glog.Info("Waiting for all controllers to shut down gracefully")
	wg.Wait()

	return nil
}
