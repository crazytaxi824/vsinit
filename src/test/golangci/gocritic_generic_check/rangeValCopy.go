package gocriticcheck

import "fmt"

// rangeValCopy
func rangeValCopy() {
	xs := make([][1024]byte, length)
	for _, x := range xs {
		fmt.Println(x)
		// Loop body.
	}
}

// panic ❌
// func rangeValCopyGeneric[T Byt]() {
// 	xs := make([][1024]T, length)
// 	for _, x := range xs {
// 		fmt.Println(x)
// 		// Loop body.
// 	}
// }
