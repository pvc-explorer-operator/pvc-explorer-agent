// Package pvcwatch monitors PVC usage status by watching pod volumes.
package pvcwatch

import (
	"context"
	"os"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Watcher monitors PVC usage by listing pods that reference the same claim.
type Watcher struct {
	client    kubernetes.Interface
	namespace string
	selfName  string
	pvcName   string
}

// New creates a Watcher that monitors PVC usage for the given claim name.
func New(pvcName string) (*Watcher, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	selfName := os.Getenv("POD_NAME")
	if selfName == "" {
		if b, err := os.ReadFile("/etc/hostname"); err == nil {
			selfName = strings.TrimSpace(string(b))
		}
	}
	return &Watcher{
		client:    client,
		namespace: os.Getenv("POD_NAMESPACE"),
		selfName:  selfName,
		pvcName:   pvcName,
	}, nil
}

// PodName returns the name of the current pod.
func (w *Watcher) PodName() string { return w.selfName }

// PVCInUse checks whether the PVC is currently mounted by another active pod.
func (w *Watcher) PVCInUse(ctx context.Context) (bool, error) {
	pods, err := w.client.CoreV1().Pods(w.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, pod := range pods.Items {
		if pod.Name == w.selfName {
			continue
		}
		if !isActive(pod.Status.Phase) {
			continue
		}
		for _, vol := range pod.Spec.Volumes {
			if vol.PersistentVolumeClaim != nil && vol.PersistentVolumeClaim.ClaimName == w.pvcName {
				return true, nil
			}
		}
	}
	return false, nil
}

func isActive(phase corev1.PodPhase) bool {
	return phase == corev1.PodRunning || phase == corev1.PodPending
}
