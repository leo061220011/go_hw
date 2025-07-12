// Задача 1: Реализация системы оплаты с использованием интерфейсов
// Описание: Создайте интерфейс PaymentProcessor с методом Process(amount float64) string,
// который возвращает строку с результатом обработки платежа. Реализуйте два типа платежных систем: CreditCard и CryptoWallet,
// каждый из которых реализует интерфейс PaymentProcessor.
// В функции main создайте список платежных систем и вызовите метод Process для каждого,
// выводя результат на экран.
// Требования:
//     Интерфейс PaymentProcessor.
//     Структуры CreditCard и CryptoWallet, реализующие интерфейс.
//     В main создайте массив/слайс этих структур и вызовите их методы.
//     P.S. Подумайте, какие поля могут быть у каждой структуры

package main

import "fmt"

type PaymentProcessor interface {
	Process(amount float64) string
}

type CreditCard struct {
	CardNumber string
	CardHolder string
	Balance    float64
	Currency   string
}

func (cc CreditCard) Process(amount float64) string {
	if amount > cc.Balance {
		return fmt.Sprintf("Недостаточно средств на карте %s", cc.CardNumber)
	}
	cc.Balance -= amount
	return fmt.Sprintf("Оплата %.2f %s, с карты %s, %s, произведена успешна.", amount, cc.Currency, cc.CardNumber, cc.CardHolder)
}

type CryptoWallet struct {
	Address    string
	WalletType string
	Balance    float64
	OwnerID    int
}

func (cw CryptoWallet) Process(amount float64) string {
	if amount > cw.Balance {
		return fmt.Sprintf("Недостаточно средств на кошельке %s", cw.Address)
	}
	cw.Balance -= amount
	return fmt.Sprintf("Оплата %.2f %s, с кошелька %s, ID владельца - %d, произведена успешна.", amount, cw.WalletType, cw.Address, cw.OwnerID)
}

func main() {

	processors := []PaymentProcessor{

		CreditCard{
			CardNumber: "2255****42552",
			CardHolder: "Ирина О.",
			Balance:    1000.00,
			Currency:   "eur",
		},

		CryptoWallet{
			Address:    "0x58C7656EC7ab88b098defB751B7401B5f6d8975P",
			WalletType: "USDT",
			Balance:    5500.00,
			OwnerID:    1234567,
		},
	}

	amounts := []float64{250.00, 50.55}
	for i, processor := range processors {
		res := processor.Process(amounts[i])
		fmt.Println(res)
	}
}
