package error_wrong_alias_single

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
)

func main() {
	k8sDeployment := appsv1.Deployment{}
	fmt.Printf("%v\n", k8sDeployment)
}
