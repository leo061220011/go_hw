// Задача 2: Обработка данных с помощью каналов
// Описание: Создайте программу, в которой одна горутина генерирует числа от 1 до 10 и отправляет их через канал.
// Другая горутина читает числа из канала и выводит их на экран. После отправки всех чисел закройте канал и завершите программу.
// Требования:

//     Используйте каналы для передачи данных.
//     Горутина-генератор должна отправлять числа и закрывать канал после этого.
//     Горутина-потребитель должна читать из канала и выводить числа, пока канал не закрыт.

package main

import (
	"fmt"
	"time"
)

func sendMessage(ch chan int) {
	for i := 1; i <= 10; i++ {
		ch <- i
	}
	close(ch)
}
func getMessage(ch chan int) {
	for msg := range ch {
		fmt.Println("Получено сообщение:", msg)
	}

}

func main() {
	ch := make(chan int)
	go sendMessage(ch)
	go getMessage(ch)
	time.Sleep(1 * time.Second)

}
