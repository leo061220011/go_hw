// Задача 2: Манипуляции со слайсами
// Создайте слайс строк, содержащий названия городов. Реализуйте функции для добавления нового города,
// удаления города по имени и поиска города в списке.
// Продемонстрируйте работу этих функций на примере.
package main

import "fmt"

func addCity(city string, citys []string) []string {

	return append(citys, city)

}
func delCity(cityToRemove string, citys []string) []string {
	for i, city := range citys {
		if city == cityToRemove {

			return append(citys[:i], citys[i+1:]...)
		}
	}
	return citys
}
func searchCity(cityToSearch string, citys []string) string {
	var c string
	for _, city := range citys {
		if city == cityToSearch {
			c = city

		}
	}
	return c
}

func main() {

	citys := []string{"Moscow", "Tula", "Omsk"}

	fmt.Println("Список городов: ")
	for _, city := range citys {
		fmt.Println(city)
	}
	citys = addCity("Perm", citys)
	fmt.Println("Добавить город Perm: ")
	for _, city := range citys {
		fmt.Println(city)
	}

	citys = delCity("Tula", citys)
	fmt.Println("Удалить город Tula: ")
	for _, city := range citys {
		fmt.Println(city)
	}
	city := searchCity("Omsk", citys)
	fmt.Println("Найти  город Omsk: ")
	fmt.Println(city)

}
