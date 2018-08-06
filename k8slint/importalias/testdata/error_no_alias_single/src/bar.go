package error_no_alias_single

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	options := v1.GetOptions{}
	fmt.Printf("%v\n", options)
}
