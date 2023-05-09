package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var serviceNamespace = "default"

func sergiceName(ld *LeptonDeployment) string {
	return ld.Name + "-service"
}

func createService(ld *LeptonDeployment, ph *Photon, or metav1.OwnerReference) error {
	clientset := mustInitK8sClientSet()

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            sergiceName(ld),
			OwnerReferences: []metav1.OwnerReference{or},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"photon": uniqName(ph.Name, ph.ID),
			},
			Ports: []corev1.ServicePort{
				{
					Port:     8080,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}

	result, err := clientset.CoreV1().Services(serviceNamespace).Create(
		context.Background(),
		service,
		metav1.CreateOptions{},
	)
	if err != nil {
		return err
	}

	fmt.Printf("Created service %q.\n", result.GetName())

	return nil
}
