package main

import (
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"context"
	"log"
	kcorev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"encoding/json"
)

var cfg = config.GetConfigOrDie()

var (
	informerCache       cache.Cache
	informerCacheCtx    context.Context
	informerCacheCancel context.CancelFunc
	knownPod1           client.Object
	knownPod2           client.Object
	knownPod3           client.Object
	knownPod4           client.Object
)

func main() {
	test1()
}


func test1() {
	ctx := context.Background()
	ns := "ns"
	knownPod1 = createPod("test-pod-1", ns, kcorev1.RestartPolicyNever)
	knownPod2 = createPod("test-pod-2", ns, kcorev1.RestartPolicyAlways)
	knownPod3 = createPod("test-pod-3", ns, kcorev1.RestartPolicyOnFailure)
	knownPod4 = createPod("test-pod-4", ns, kcorev1.RestartPolicyNever)
	podGVK := schema.GroupVersionKind{
		Kind:    "Pod",
		Version: "v1",
	}
	knownPod1.GetObjectKind().SetGroupVersionKind(podGVK)
	knownPod2.GetObjectKind().SetGroupVersionKind(podGVK)
	knownPod3.GetObjectKind().SetGroupVersionKind(podGVK)
	knownPod4.GetObjectKind().SetGroupVersionKind(podGVK)


	informerCache, _ := cache.New(cfg, cache.Options{})
	go informerCache.Start(ctx)

	err := informerCache.WaitForCacheSync(ctx)
	log.Printf("wait for sync with error : %v\n", err )

	outPod := &kcorev1.Pod{}

	err = informerCache.Get(ctx, "ns/test-pod-3", outPod)
	pretty_pod, _ := json.MarshalIndent(outPod, "", "\t")
	log.Printf("pretty_pod: %v\n err: %v\n", string(pretty_pod), err)
}

func createPod(name, namespace string, restartPolicy kcorev1.RestartPolicy) client.Object {
	three := int64(3)
	pod := &kcorev1.Pod{
		ObjectMeta: kmetav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"test-label": name,
			},
		},
		Spec: kcorev1.PodSpec{
			Containers:            []kcorev1.Container{{Name: "nginx", Image: "nginx"}},
			RestartPolicy:         restartPolicy,
			ActiveDeadlineSeconds: &three,
		},
	}
	cl, err := client.New(cfg, client.Options{})
	err = cl.Create(context.Background(), pod)
	log.Printf("createPod with err: %v\n", err)
	return pod
}


func deletePod(pod client.Object) {
	cl, err := client.New(cfg, client.Options{})
	err = cl.Delete(context.Background(), pod)
	log.Printf("deletePod with err: %v\n", err)
}

