// Задача 3: Структура "Банковский счет" и методы для работы с балансом
// Описание:
// Создайте структуру BankAccount с полями: Owner (владелец счета) и Balance (текущий баланс).
// Реализуйте методы:
// Deposit(amount float64), увеличивающий баланс.
// Withdraw(amount float64), уменьшающий баланс, если хватает средств, иначе выводит сообщение о недостатке средств.
// Что нужно сделать:
// Объявить структуру и методы.
// Создать счет, пополнить его, попытаться снять деньги и вывести итоговый баланс.
package main

import "fmt"

type BankAccount struct {
	Owner   string
	Balance float64
}

func (b *BankAccount) Deposit(amount float64) {
	b.Balance += amount
}

func (b *BankAccount) Withdraw(amount float64) {
	if b.Balance >= amount {
		b.Balance -= amount
	} else {
		fmt.Println("Отказ! Средств недостаточно!")
	}

}

func main() {

	b := BankAccount{Owner: "Александр", Balance: 1000.00}
	fmt.Printf("Владелец счета: %s\nТекущий баланс: %v\n", b.Owner, b.Balance)
	b.Deposit(100.00)
	fmt.Printf("Пополнение счета: 100.00\nТекущий баланс: %v\n", b.Balance)
	fmt.Println("Снятие со счета: 2000.00")
	b.Withdraw(2000.00)
	fmt.Printf("Текущий баланс: %v\n", b.Balance)
	fmt.Println("Снятие со счета: 200.00")
	b.Withdraw(200.00)
	fmt.Printf("Текущий баланс: %v\n", b.Balance)

}
