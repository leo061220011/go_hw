// Задача 3
// Напишите программу, которая:

// 1) Создаёт одну анонимную функцию, которая при каждом вызове увеличивает внутренний счётчик на 1 и возвращает его значение.
// 2) Вызывает эту функцию несколько раз и выводит результат.
// Пример:
//
//	// Вызовы функции
//	fmt.Println(counter()) // 1
//	fmt.Println(counter()) // 2
//	fmt.Println(counter()) // 3
package main

import "fmt"

func main() {
	counter := func() func() int {
		count := 0
		return func() int {
			count++
			return count
		}
	}()

	fmt.Println(counter())
	fmt.Println(counter())
	fmt.Println(counter())
}
