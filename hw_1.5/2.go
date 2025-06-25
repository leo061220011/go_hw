// // Задача 2:
// // Напишите программу, которая подсчитывает количество гласных букв (а, е, ё, и, о, у, ы, э, ю, я) в введённой пользователем строке.

package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func main() {

	var str string

	fmt.Println("Введите строку: ")
	str, _ = bufio.NewReader(os.Stdin).ReadString('\n')

	glasnye := []rune{'а', 'е', 'ё', 'и', 'о', 'у', 'ы', 'э', 'ю', 'я'}

	count := 0
	for _, element := range str {
		change := unicode.ToLower(element)
		for _, element1 := range glasnye {
			if change == element1 {
				count++
			}

		}

	}
	fmt.Printf("Количество гласных букв: %v", count)

}

//package main

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// )

// func main() {
// 	text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
// 	fmt.Println(text)
// }
