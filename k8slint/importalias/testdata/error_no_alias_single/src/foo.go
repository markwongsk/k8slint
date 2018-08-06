package error_no_alias_single

import (
	"fmt"

	"k8s.io/api/apps/v1"
)

func main() {
	k8sDeployment := v1.Deployment{}
	fmt.Printf("%v\n", k8sDeployment)
}
