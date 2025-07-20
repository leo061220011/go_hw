// Задача 1
// Создайте программу, которая читает лог-файл server.log,
// подсчитывает количество строк, содержащих слово "error",
// и выводит это число.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

func main() {

	file, err := os.Open("hw_1.12/server.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	re := regexp.MustCompile(`\berror\b`)

	scanner := bufio.NewScanner(file)
	count := 0

	for scanner.Scan() {

		line := scanner.Text()
		if re.MatchString(line) {
			count++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Количество строк, содержащих слово `error`: %d\n", count)
}
