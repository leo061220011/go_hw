// Задача 1
// Реализуйте функцию divide(a, b float64) (float64, error), которая делит a на b.
//
//	Если b равно нулю, возвращайте ошибку.
//	В основной программе вызовите эту функцию и обработайте возможную ошибку.
package main

import "fmt"

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("деление на ноль")
	}
	return a / b, nil
}

func main() {

	result, err := divide(155, 0)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Результат:", result)
	}

}
