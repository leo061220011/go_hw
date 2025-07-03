// Задача 3: Использование мапов для подсчета частот
// Напишите программу, которая читает строку текста и подсчитывает количество вхождений каждого слова.
// Используйте мапу (map[string]int) для хранения результатов.
// Выведите полученную статистику.
package main

import (
	"fmt"
	"strings"
)

func main() {

	var text string = "Сара - Абрамовна - дюдюка"
	words := strings.Fields(text)
	res := make(map[string]int)
	for _, word := range words {
		res[word]++
	}
	for key, value := range res {
		if key != "-" {
			fmt.Printf("%s %d\n", key, value)
		}
	}

}
