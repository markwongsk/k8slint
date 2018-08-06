package no_error_multiple

import (
	"fmt"

	apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	k8sDeployment := apps.Deployment{}
	job := batch.Job{}
	option := meta.GetOptions{}
	fmt.Printf("%v, %v, %v\n", k8sDeployment, job, option)
}
