package cmd

import (
	awsv1alpha1 "github.com/openshift/aws-account-operator/pkg/apis/aws/v1alpha1"

	"k8s.io/client-go/kubernetes/scheme"
)

func init() {
	awsv1alpha1.AddToScheme(scheme.Scheme)
}
