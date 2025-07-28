// Задача 4 из темы 8
// Создайте функцию sumAll, которая принимает произвольное количество целых чисел и возвращает их сумму.
// Пример использования:
//     fmt.Println(sumAll(1, 2, 3)) // 6
//     fmt.Println(sumAll(10, -2, 4, 7)) // 19

package main

import "testing"

func TestSumAll(t *testing.T) {
	if res1 := sumAll(1, 2, 3); res1 != 6 {
		t.Errorf("sumAll(1, 2, 3) = %d, должно 6", res1)
	}
	if res2 := sumAll(10, -2, 4, 7); res2 != 19 {
		t.Errorf("sumAll(10, -2, 4, 7) = %d, ljk;yj 19", res2)
	}
}
