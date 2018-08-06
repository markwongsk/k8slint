package no_error_single

import (
	"fmt"

	apps "k8s.io/api/apps/v1"
)

func main() {
	k8sDeployment := apps.Deployment{}
	fmt.Printf("%v\n", k8sDeployment)
}
