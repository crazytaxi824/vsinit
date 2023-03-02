package gocriticcheck

import "fmt"

// rangeExprCopy
func rangeExprCopy() {
	var xs [2048]byte
	for _, x := range xs { // Copies 2048 bytes
		fmt.Println(x)
		// Loop body.
	}
}

// panic ❌
// func rangeExprCopyGeneric[T Byt]() {
// 	var xs [2048]T
// 	for _, x := range xs { // Copies 2048 bytes
// 		fmt.Println(x)
// 		// Loop body.
// 	}
// }
