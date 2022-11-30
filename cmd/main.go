package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

type BankClient interface {
	// Deposit Внесите заданную сумму на счет клиента
	Deposit(amount int)

	// Withdrawal Вывод средств выводит указанную сумму со счета клиента.
	// Ошибка возврата, если баланс клиента меньше суммы вывода
	Withdrawal(amount int) error

	// Balance Баланс возвращает баланс клиентов
	Balance() int
}

func NewBankClient() BankClient {
	return &bank{}
}

type bank struct {
	mx    sync.RWMutex
	value int
}

func (b *bank) Deposit(amount int) {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.value += amount // b.value = v.value + amount
}

func (b *bank) Withdrawal(amount int) error {
	b.mx.Lock()
	defer b.mx.Unlock()

	if b.value < amount {
		return errors.New("баланс не достаточен для снятия, операция не может быть выполнена")
	}

	b.value -= amount

	return nil
}

func (b *bank) Balance() int {
	return b.value
}

func main() {
	maxValue := 100000
	rand.Seed(time.Now().UTC().UnixNano())
	g := NewBankClient()

	wgDeposit, wgWithdrawal := sync.WaitGroup{}, sync.WaitGroup{}
	for i := 1; i <= 10; i++ {
		wgDeposit.Add(1)
		go func() {
			n := rand.Intn(maxValue)
			g.Deposit(n)
			wgDeposit.Done()
		}()
	}

	for j := 1; j <= 5; j++ {
		wgWithdrawal.Add(1)
		go func() {
			n := rand.Intn(maxValue)
			err := g.Withdrawal(n)
			if err != nil {
				fmt.Println(err)
			}
			wgWithdrawal.Done()
		}()
	}

	wgDeposit.Wait()
	wgWithdrawal.Wait()

	fmt.Println("-----------------------------")

	fmt.Println("Вы можете использовать команды: balance , deposit , withdrawal , exit")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		command := scanner.Text()

		switch BankCommands(command) {
		case balanceCommand:
			fmt.Printf("Текущий Баланс: %d \n", g.Balance())
		case exitCommand:
			fmt.Println("программа завершила работу")
			os.Exit(1) // Типав ctrl/c
		case depositCommand:
			fmt.Println("Введите сумму зачисления")
			amount, err := scanAmount(scanner)
			if err != nil {
				fmt.Println(err)
				fmt.Println("введите команду заново")
				continue
			}

			g.Deposit(amount)
			fmt.Printf("Текущий Баланс: %d \n", g.Balance())
		case withdrawalCommand:
			fmt.Println("Введите сумму снятия")
			amount, err := scanAmount(scanner)
			if err != nil {
				fmt.Println(err)
				fmt.Println("введите команду заново")
				continue
			}

			if err := g.Withdrawal(amount); err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Текущий Баланс: %d \n", g.Balance())
		default:
			fmt.Println("неизвестная команда")
			fmt.Println("Вы можете использовать команды: balance , deposit , withdrawal , exit")
			continue
		}
	}
}

func scanAmount(scanner *bufio.Scanner) (int, error) {
	scanner.Scan()
	amount, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return 0, errors.New("некорректная сумма")
	}

	return amount, nil
}

type BankCommands string

const (
	balanceCommand    BankCommands = "balance"
	depositCommand    BankCommands = "deposit"
	withdrawalCommand BankCommands = "withdrawal"
	exitCommand       BankCommands = "exit"
)
