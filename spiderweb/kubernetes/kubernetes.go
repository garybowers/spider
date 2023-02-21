package kubernetes

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func clusterClient() *kubernetes.Clientset {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		log.Println(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return clientset
}

func GetDeployments(namespace string) *appsv1.DeploymentList {
	deploymentsClient := clusterClient().AppsV1().Deployments(namespace)
	fmt.Printf("Listing deployments in namespace %q:\n", namespace)
	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
	}
	return list
}

func GetServices(namespace string) *apiv1.ServiceList {
	svcClient := clusterClient().CoreV1().Services(namespace)
	list, err := svcClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
	}
	return list
}

func CreatePersistentVolume(namespace string, pvspec *apiv1.PersistentVolume) {
	result, err := clusterClient().CoreV1().PersistentVolumes().Create(context.TODO(), pvspec, metav1.CreateOptions{})
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("Created persistent volume %q.\n", result.GetObjectMeta().GetName())
	return
}

func CreatePersistentVolumeClaim(namespace string, pvcspec *apiv1.PersistentVolumeClaim) {
	result, err := clusterClient().CoreV1().PersistentVolumeClaims(namespace).Create(context.TODO(), pvcspec, metav1.CreateOptions{})
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("Created persistent claim %q.\n", result.GetObjectMeta().GetName())
	return
}

func CreateService(namespace string, svcspec *apiv1.Service) {
	result, err := clusterClient().CoreV1().Services(namespace).Create(context.TODO(), svcspec, metav1.CreateOptions{})
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("Created service %q.\n", result.GetObjectMeta().GetName())
	return
}

func CreateDeployment(namespace string, deploymentspec *appsv1.Deployment) {
	result, err := clusterClient().AppsV1().Deployments(namespace).Create(context.TODO(), deploymentspec, metav1.CreateOptions{})
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Creating deployment %s.\n", result.GetObjectMeta().GetName())
	d, err := clusterClient().AppsV1().Deployments(namespace).Get(context.TODO(), result.GetObjectMeta().GetName(), metav1.GetOptions{})
	if err != nil {
		log.Println(err.Error())
	}

	var i int = 0
	for d.Status.Replicas == 0 {
		if i == 1000 {
			log.Println("Timeout waiting for deployment %s.", d.Name)
			break
		}
		fmt.Print("Deployment: %s is waiting to become available....", string(d.Name))
		d, err = clusterClient().AppsV1().Deployments(namespace).Get(context.TODO(), result.GetObjectMeta().GetName(), metav1.GetOptions{})
		if err != nil {
			log.Println(err.Error())
		}
	}

	log.Printf("Finished Creating deployment %s.\n", result.GetObjectMeta().GetName())
	return
}

func DeleteDeployment(namespace string, deploymentName string) {
	deletePolicy := metav1.DeletePropagationForeground
	if err := clusterClient().AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Printf("Error deleting deployment %s.\n", err)
	}
	log.Printf("Deleted deployment %s.\n", deploymentName)
	return
}

func DeleteService(namespace string, serviceName string) {
	deletePolicy := metav1.DeletePropagationForeground
	if err := clusterClient().CoreV1().Services(namespace).Delete(context.TODO(), serviceName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Printf("Error deleting service %s.\n", err)
	}
	log.Printf("Deleted service %s.\n", serviceName)
	return
}

func DeletePersistentVolumeClaim(namespace string, pvcName string) {
	deletePolicy := metav1.DeletePropagationForeground
	if err := clusterClient().CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), pvcName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	log.Printf("Deleted PersistentVolumeClaim %s.\n", pvcName)
	return
}
