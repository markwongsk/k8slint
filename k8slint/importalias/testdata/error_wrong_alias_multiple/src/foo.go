package error_wrong_alias_multiple

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	k8sDeployment := appsv1.Deployment{}
	job := batchv1.Job{}
	options := metav1.GetOptions{}
	fmt.Printf("%v, %v, %v\n", k8sDeployment, job, options)
}
