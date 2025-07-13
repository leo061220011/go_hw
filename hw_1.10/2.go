// Задача 2: Полиморфное отображение данных
// Описание: Создайте интерфейс Shape с методом Area() float64.
// Реализуйте структуры Circle и Rectangle, которые реализуют этот интерфейс.
// Напишите функцию, которая принимает слайс фигур ([]Shape) и выводит площадь каждой фигуры.
// Требования:
// Интерфейс Shape.
// Структуры Circle (с радиусом) и Rectangle (с длиной и шириной).
// Функция для вывода площадей всех фигур в слайсе.
package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() (float64, string)
}

type Circle struct {
	Radius float64
	Name   string
}

func (c Circle) Area() (float64, string) {
	a := math.Pi * c.Radius * c.Radius
	n := c.Name
	return a, n
}

type Rectangle struct {
	Width, Height float64
	Name          string
}

func (r Rectangle) Area() (float64, string) {
	a := r.Width * r.Height
	n := r.Name
	return a, n
}

func myArea(s Shape) {
	a, n := s.Area()
	fmt.Printf("Площадь фигуры %s: %2.f\n", n, a)
}

func main() {
	Shapes := []Shape{
		Circle{Name: "круг", Radius: 10},
		Rectangle{Name: "прямоугольник", Width: 4, Height: 7},
	}
	for _, shape := range Shapes {
		myArea(shape)
	}
}
