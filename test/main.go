/*
 * @time 2021/10/10 7:43 下午
 * @version 1.00
 * @author huangsiyi
 */
package main

import (
	"fmt"
)

func main() {
	a := []int{}
	b := []int{1, 2, 3}
	a = b
	// c: nil
	// a = append(b, 1) // a: [1, 2, 3, 1]
	c := a
	a = append(a, 1)

	fmt.Println("a: ", a)
	fmt.Println("b: ", b)
	fmt.Println("c: ", c)
	// 修改后，如果b达到空间最大值，分配新空间后，c仍指向原来a的地址
}
