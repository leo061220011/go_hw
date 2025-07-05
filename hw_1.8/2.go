// Задача 2
// Функция высшего порядка: передача функции как аргумента
// Создайте функцию applyOperation(a, b int, op func(int, int) int) int,
//  которая применяет переданную функцию op к числам a и b.

// Создайте несколько функций-операций: сложение, вычитание, умножение.
// В основной программе вызовите applyOperation с разными операциями и выведите результаты.
package main

import "fmt"

func applyOperation(a, b int, op func(int, int) int) int {
	return op(a, b)
}

func add(a, b int) int {
	return a + b
}

func multiply(a, b int) int {
	return a * b
}

func subtract(a, b int) int {
	return a - b
}
func main() {
	res := applyOperation(5, 10, add)
	fmt.Printf("Результат сложения: %d\n", res)

	res = applyOperation(5, 10, multiply)
	fmt.Printf("Результат умножения: %d\n", res)

	res = applyOperation(5, 10, subtract)
	fmt.Printf("Результат вычитания: %d\n", res)
}
