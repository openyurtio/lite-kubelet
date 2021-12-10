package client

import (
	"context"
	"path/filepath"

	"k8s.io/klog/v2"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type PodsGetter interface {
	Pods(namespace string) PodInstance
}

type PodInstance interface {
	PublishTopicor
	Create(ctx context.Context, pod *corev1.Pod, opts v1.CreateOptions) (result *corev1.Pod, err error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *corev1.Pod, err error)
}

type pods struct {
	nodename  string
	namespace string
	client    MessageSendor
}

func (p *pods) GetPublishDeleteTopic(name string) string {
	return filepath.Join(p.GetPublishPreTopic(), name, "delete")
}

func (p *pods) GetPublishCreateTopic(name string) string {
	return filepath.Join(p.GetPublishPreTopic(), name, "create")
}

func (p *pods) GetPublishUpdateTopic(name string) string {
	return filepath.Join(p.GetPublishPreTopic(), name, "update")
}

func (p *pods) GetPublishPatchTopic(name string) string {
	return filepath.Join(p.GetPublishPreTopic(), name, "patch")
}

func (p *pods) GetPublishPreTopic() string {
	return filepath.Join(MqttEdgePublishRootTopic, "pods", p.namespace)
}

func (p *pods) Create(ctx context.Context, pod *corev1.Pod, opts v1.CreateOptions) (result *corev1.Pod, err error) {
	klog.Warningf("implement me, create pod")
	return pod, nil
}

func (p *pods) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	klog.Warningf("implement me, delete pod")
	return nil
}

func (p *pods) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *corev1.Pod, err error) {
	klog.Warningf("implement me, patch pod")
	return nil, nil
}

func newPods(nodename, namespace string, c MessageSendor) *pods {
	return &pods{
		nodename:  nodename,
		namespace: namespace,
		client:    c,
	}
}

var _ PodInstance = &pods{}
