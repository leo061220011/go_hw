package main

import "fmt"

func main() {

	var a, b, c, d, e int

	a, b, c, d, e = 55, 4, 75, 3, 15

	for i := 0; i < 5; i++ {

		if a < b {
			a, b = b, a
		}
		if b < c {
			b, c = c, b
		}
		if c < d {
			c, d = d, c
		}
		if d < e {
			d, e = e, d
		}

	}

	var z int = (a + b + c + d + e) / 5

	fmt.Printf("Отсортированые элементы: %d %d %d %d %d\n", a, b, c, d, e)

	fmt.Printf("Самое большое число: %d\n", a)

	fmt.Printf("Самое маленькое число: %d\n", e)

	fmt.Printf("Среднее-арифметическое: %d\n", z)

}
