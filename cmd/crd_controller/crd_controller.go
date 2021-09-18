package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//start app

	s.AddFlags(pflag.CommandLine)

	klog.InitFlags(nil)
	flag.Parse()

	stopCh := signals.SetupSignalHandler()

	verflag.PrintAndExitIfRequested()

	if err := s.Run(pflag.CommandLine.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
