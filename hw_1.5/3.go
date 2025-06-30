// Задача 3:
// Создайте функцию capitalizeWords(s string) string, которая преобразует каждое слово в строке так, чтобы первая буква была заглавной, а остальные — строчными. Например: "привет мир" → "Привет Мир".
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func capitalizeWords(s string) string {

	str := strings.ToLower(s)
	words := strings.Fields(str)

	for i, word := range words {
		runes := []rune(word)
		if len(runes) > 0 {

			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}

	return strings.Join(words, " ")
}

func main() {
	var stroca string
	stroca, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Println(capitalizeWords(stroca))

}
