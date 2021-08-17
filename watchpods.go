package main
import (
    "flag"
    "fmt"
    "time"
    "k8s.io/api/core/v1"
    "k8s.io/client-go/kubernetes"
	"k8s.io/apimachinery/pkg/fields"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/tools/clientcmd"
)
var (
    kubeconfig = flag.String("kubeconfig", "khushi-watch-test", "absolute path to the kubeconfig file")
)
func main() {
    flag.Parse()

	// Build a working client from the kubeconfig
    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
        panic(err.Error())
    }

	// Create a new clientset for the given config
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err.Error())
    }


    watchlist := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", v1.NamespaceDefault, fields.Everything())
    _, controller := cache.NewInformer(
        watchlist,
        &v1.Pod{},
        time.Second*0,
        cache.ResourceEventHandlerFuncs{
            AddFunc: func(obj interface{}) {
                fmt.Printf("pod added: %s \n", obj)
            },
            DeleteFunc: func(obj interface{}) {
                fmt.Printf("pod deleted: %s \n", obj)
            },
            UpdateFunc: func(oldObj, newObj interface{}) {
                fmt.Printf("pod changed \n")
            },
        },
    )
    stop := make(chan struct{})
    go controller.Run(stop)
    for {
        time.Sleep(time.Second)
    }
}
