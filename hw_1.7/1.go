// Задача 1: Работа с массивами и слайсами
// Напишите программу, которая создает массив из 10 целых чисел, заполняет его случайными значениями от 1 до 100.
// Затем скопируйте этот массив в слайс и отсортируйте его по возрастанию. Выведите исходный массив и отсортированный слайс.

package main

import (
	"fmt"
	"math/rand"
	"sort"
)

func main() {

	var numbers [10]int
	for i := range numbers {

		numbers[i] = rand.Intn(101)
	}
	numbers_slice := make([]int, len(numbers))
	copy(numbers_slice, numbers[:])
	sort.Ints(numbers_slice)
	fmt.Println("Исходный массив: ")
	for _, value := range numbers {
		fmt.Println(value)
	}
	fmt.Println("Отсортированный слайс: ")
	for _, value := range numbers_slice {
		fmt.Println(value)
	}
}
