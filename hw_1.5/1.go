// Задача 1:
// Напишите программу, которая запрашивает у пользователя ввод строки, а затем выводит число - количество символов в строке

package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {

	var str string

	fmt.Println("Введите строку: ")
	fmt.Scanln(&str)

	res := utf8.RuneCountInString(str)

	fmt.Println(res)

}
