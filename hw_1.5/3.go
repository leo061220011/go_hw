// Задача 3:
// Создайте функцию capitalizeWords(s string) string, которая преобразует каждое слово в строке так, чтобы первая буква была заглавной, а остальные — строчными. Например: "привет мир" → "Привет Мир".
package main

import (
	"bufio"
	"fmt"
	"os"
)

var capitalizeWords(s string) string {
w := strings.Fields(s)
for i, b := range w {
	if len(b) == 0 {
	continue
}
}
}

func main() {

	var str string


	fmt.Println("Введите строку: ")
	str, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Println(CapitalizeWords("str"))


}
