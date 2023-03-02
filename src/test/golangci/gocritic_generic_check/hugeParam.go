package gocriticcheck

import "fmt"

// hugeParam
func hugeParam(x [1024]int) {
	fmt.Println(x)
}

// panic ❌
// func hugeParamGeneric[T Ints](x [1024]T) {
// 	fmt.Println(x)
// }
