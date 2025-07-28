package main

import "fmt"

func sumAll(nums ...int) int {
	res := 0
	for _, num := range nums {
		res += num
	}
	return res
}

func main() {

	fmt.Println(sumAll(1, 2, 3))
	fmt.Println(sumAll(10, -2, 4, 7))

}
