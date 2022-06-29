package main

import (
	"context"
	"fmt"
	"io"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"os"
)

type KubeClient struct {
	client     kubernetes.Interface
	restConfig *rest.Config
}

func (k *KubeClient) CopyByLabel(label string, namespace string, container string, containerPath string, useFirst bool) (*io.PipeReader, error) {

	pods, err := k.client.CoreV1().Pods(namespace).List(context.Background(), v12.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, err
	}
	fmt.Printf("Pod count:%v\n", len(pods.Items))
	if len(pods.Items) > 1 {
		if useFirst {
			fmt.Printf("multiple pods found: %v, using Pod: %v \n", len(pods.Items), pods.Items[0].Name)
		} else {
			return nil, fmt.Errorf("only 1 pod expected")
		}
	}
	if len(pods.Items) == 0 {
		return nil, fmt.Errorf("no pods found")
	}
	return k.CopyByPodName(pods.Items[0].Name, namespace, container, containerPath)
}

func (k *KubeClient) CopyByPodName(podName string, namespace string, container string, containerPath string) (*io.PipeReader, error) {
	req := k.client.CoreV1().RESTClient().Post().Resource("pods").Name(podName).
		Namespace(namespace).SubResource("exec")
	return k.copyFromPod(req, container, containerPath)
}

func (k *KubeClient) copyFromPod(req *rest.Request, container string, containerPath string) (*io.PipeReader, error) {
	fmt.Printf("copying from pod: %v\n", req.URL())
	reader, outStream := io.Pipe()
	cmd := []string{"tar", "cf", "-", containerPath}
	option := &v1.PodExecOptions{
		Command:   cmd,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		Container: container,
		TTY:       false,
	}
	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	exec, err := remotecommand.NewSPDYExecutor(k.restConfig, "POST", req.URL())
	if err != nil {
		return nil, err
	}
	go func() {
		defer func(outStream *io.PipeWriter) {
			err := outStream.Close()
			if err != nil {
				fmt.Printf("error closing file: %v\n", err)
			}
		}(outStream)
		err = exec.Stream(remotecommand.StreamOptions{
			Stdin:  os.Stdin,
			Stdout: outStream,
			Stderr: os.Stderr,
		})
		if err != nil {
			fmt.Printf("Error calling kubernetes %v\n", err)
		}
	}()
	return reader, nil
}

func (k *KubeClient) CopyFromPod(podName string, namespace string, container string, containerPath string) (*io.PipeReader, error) {
	fmt.Printf("Copy from POD\n")

	req := k.client.CoreV1().RESTClient().Post().Resource("pods").Name(podName).
		Namespace(namespace).SubResource("exec")
	reader, outStream := io.Pipe()
	cmd := []string{"tar", "cf", "-", containerPath}
	option := &v1.PodExecOptions{
		Command:   cmd,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		Container: container,
		TTY:       false,
	}
	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	exec, err := remotecommand.NewSPDYExecutor(k.restConfig, "POST", req.URL())
	if err != nil {
		return nil, err
	}
	go func() {
		defer func(outStream *io.PipeWriter) {
			err := outStream.Close()
			if err != nil {
				fmt.Printf("error closing file: %v\n", err)
			}
		}(outStream)
		err = exec.Stream(remotecommand.StreamOptions{
			Stdin:  os.Stdin,
			Stdout: outStream,
			Stderr: os.Stderr,
		})
		if err != nil {
			fmt.Printf("Error calling kubernetes %v\n", err)
		}
	}()
	return reader, nil
}
