package main

import (
	"flag"
	"github.com/nistal97/crd_controller/internal"
	clientset "github.com/nistal97/crd_controller/pkg/generated/clientset/versioned"
	informers "github.com/nistal97/crd_controller/pkg/generated/informers/externalversions"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"runtime"
	"time"

	"github.com/nistal97/crd_controller/pkg/signals"
	"k8s.io/klog/v2"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	//uncomment me when debugging
	//masterURL := "https://api.system.svc.130.tess.io"
	//kubeconfig := "./kubeconfig"

	runtime.GOMAXPROCS(runtime.NumCPU())

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err)
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err)
	}

	ciconfigClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building example clientset: %s", err)
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	ciconfigInformerFactory := informers.NewSharedInformerFactory(ciconfigClient, time.Second*30)

	controller := internal.NewCiConfigController(kubeClient, ciconfigClient,
		ciconfigInformerFactory.Tess().V1().CiConfigs())

	// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	kubeInformerFactory.Start(stopCh)
	ciconfigInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err)
	}
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
