package internal

import (
	"fmt"
	"time"

	v1 "github.com/nistal97/crd_controller/pkg/api/tess.io/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	clientset "github.com/nistal97/crd_controller/pkg/generated/clientset/versioned"
	scheme "github.com/nistal97/crd_controller/pkg/generated/clientset/versioned/scheme"
	informers "github.com/nistal97/crd_controller/pkg/generated/informers/externalversions/tess.io/v1"
	listers "github.com/nistal97/crd_controller/pkg/generated/listers/tess.io/v1"
)

const (
	controllerAgentName = "ciconfig-controller"

	// SuccessSynced is used as part of the Event 'reason' when a ciconfig is synced
	SuccessSynced = "Synced"
	// MessageResourceSynced is the message used for an Event fired when a ciconfig
	// is synced successfully
	MessageResourceSynced = "CiConfig synced successfully"
)

type CiConfigController struct {
	// standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// ciconfig clientset
	ciconfig_clientset clientset.Interface

	ciconfigLister listers.CiConfigLister
	ciconfigSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController creates a new controller
func NewCiConfigController(
	kubeclientset kubernetes.Interface,
	ciconfigclientset clientset.Interface,
	ciconfigInformer informers.CiConfigInformer) *CiConfigController {

	utilruntime.Must(scheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster..")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("ciaas-test")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &CiConfigController{
		kubeclientset:      kubeclientset,
		ciconfig_clientset: ciconfigclientset,
		ciconfigLister:     ciconfigInformer.Lister(),
		ciconfigSynced:     ciconfigInformer.Informer().HasSynced,
		workqueue:          workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ciconfig"),
		recorder:           recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when ciconfig resources change
	ciconfigInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueCiconfig,
		UpdateFunc: func(old, new interface{}) {
			oldMydemo := old.(*v1.CiConfig)
			newMydemo := new.(*v1.CiConfig)
			if oldMydemo.ResourceVersion == newMydemo.ResourceVersion {
				return
			}
			controller.enqueueCiconfig(new)
		},
		DeleteFunc: controller.enqueueCiconfigForDelete,
	})
	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *CiConfigController) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting ciconfig controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.ciconfigSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch workers to process resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers..")
	<-stopCh
	klog.Info("Shutting down workers..")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *CiConfigController) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *CiConfigController) processNextWorkItem() bool {

	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Ciconfig resource to be synced.
		if err := c.syncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the ciconfig resource
// with the current status of the resource.
func (c *CiConfigController) syncHandler(key string) error {

	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the ciconfig resource with this namespace/name
	ciconfig, err := c.ciconfigLister.CiConfigs(namespace).Get(name)
	if err != nil {
		//not found, means already deleted
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("ciconfig '%s' in work queue no longer exists", key))
			return nil
		}
		runtime.HandleError(fmt.Errorf("failed to list ciconfig by: %s/%s", namespace, name))
		return err
	}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	c.recorder.Event(ciconfig, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

// first in cache, then enqueue
// Ciconfig -> namespace/name, then put into workqueue
func (c *CiConfigController) enqueueCiconfig(obj interface{}) {
	var key string
	var err error
	// put obj in cache
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}

	klog.Infof("enqueueCiconfig:" + key)
	//put key in queue
	c.workqueue.AddRateLimited(key)
}

func (c *CiConfigController) enqueueCiconfigForDelete(obj interface{}) {
	var key string
	var err error
	// delete from cache
	key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	//put key in queue
	klog.Infof("enqueueCiconfigForDelete:" + key)
	c.workqueue.AddRateLimited(key)
}
