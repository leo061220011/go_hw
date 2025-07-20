// Задача 2
// Напишите программу, которая принимает список имен файлов в текущей директории,
// объединяет их содержимое и сохраняет результат в новый файл combined.txt (работать только с текстовыми файлами)

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	files, err := filepath.Glob("hw_1.12/newDirectory/*.txt")
	if err != nil {
		log.Fatal(err)
	}

	outputFile, err := os.Create("hw_1.12/combined.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			line := strings.TrimRight(scanner.Text(), " \t")
			if line != "" {
				_, err := fmt.Fprintf(outputFile, "%s\n", line)
				if err != nil {
					log.Printf("Ошибка записи: %v", err)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Ошибка при чтении %s: %v", file, err)
		}
		f.Close()
	}

}
