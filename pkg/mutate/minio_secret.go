package mutate

// Read data from the MinIO secret currently deployed on the cluster.

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// MinIOSecretData contains the data from the MinIO secret
type MinIOSecretData struct {
	// TODO: support encoding (base64)
	UserName     string
	UserPassword string
}

// ReadMinIOSecret reads the data from the MinIO secret and returns it as a
// MinIOSecretData
func ReadMinIOSecret() (*MinIOSecretData, error) {
	// Use the in-cluster config to create the clientSet
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get the secret by name and nameSpace
	type secretArgs struct {
		nameSpace string
		name      string
	}
	sa := secretArgs{nameSpace: "minio", name: "minio"} // yes, hardcoded now
	secret, err := clientSet.CoreV1().
		Secrets(sa.nameSpace).
		Get(ctx, sa.name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// Extract the data from the secret
	username, ok := secret.Data["root-user"]
	if !ok {
		return nil, fmt.Errorf("username key not found in secret")
	}
	password, ok := secret.Data["root-password"]
	if !ok {
		return nil, fmt.Errorf("password key not found in secret")
	}
	logrus.WithFields(logrus.Fields{
		"username": string(username),
		// "password": string(password),
	}).Info("ReadMinIOSecret()")

	return &MinIOSecretData{
		UserName:     string(username),
		UserPassword: string(password),
	}, nil
}
