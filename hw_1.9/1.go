// Задача 1: Создание структуры "Книга"
// Описание:
// Создайте структуру Book, которая содержит поля: Title (строка), Author (строка), Year (целое число).
// Добавьте метод GetInfo(), который возвращает строку с информацией о книге в формате: "Title" by Author, Year.
// Что нужно сделать:
// Объявить структуру Book.
// Реализовать метод GetInfo() для этой структуры.
// Создать экземпляр книги и вывести информацию с помощью метода.
package main

import "fmt"

type Book struct {
	Title  string
	Author string
	Year   int
}

func (b Book) GetInfo() {
	fmt.Printf("Title: %s\nBy Author: %s\nYear: %d\n", b.Title, b.Author, b.Year)

}

func main() {

	b := Book{Title: "Война и мир", Author: "Толстой", Year: 1869}

	b.GetInfo()
}
