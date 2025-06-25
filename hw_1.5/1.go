// Задача 1:
// Напишите программу, которая запрашивает у пользователя ввод строки, а затем выводит число - количество символов в строке

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {

	var str string

	fmt.Println("Введите строку: ")
	str, _ = bufio.NewReader(os.Stdin).ReadString('\n')

	str = strings.TrimRight(str, "\r\n")

	res := utf8.RuneCountInString(str)

	fmt.Println(res)

}
