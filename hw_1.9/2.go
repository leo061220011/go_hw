// Задача 2: Структура "Студент" и метод для вычисления среднего балла
// Описание:
// Создайте структуру Student, которая содержит поля: имя (Name) и список оценок (Grades []float64).
// Реализуйте метод AverageGrade() float64, который возвращает средний балл студента.
// Что нужно сделать:
// Объявить структуру и метод.
// Создать студента с несколькими оценками и вывести его средний балл.
package main

import "fmt"

type Student struct {
	Name   string
	Grades []float64
}

func (s *Student) AverageGrade() float64 {
	sum := 0.0
	for _, grade := range s.Grades {
		sum += grade
	}

	return sum / float64(len(s.Grades))

}

func main() {

	s := Student{Name: "Александр", Grades: []float64{1, 2, 3, 4, 5}}
	res := s.AverageGrade()
	fmt.Printf("Name: %s\naveragescore: %.2f", s.Name, res)
}
