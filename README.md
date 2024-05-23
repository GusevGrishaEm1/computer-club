# Компьютерный клуб - Инструкция по запуску
> Этот проект реализует сервис для управления компьютерным клубом, включая регистрацию клиентов, распределение столов,
> обработку очередей и расчет стоимости использования.
> В этом руководстве описаны шаги для установки, настройки и запуска сервиса.

# Требования
* Go 1.18 или выше
* Установленный git
* Доступ к командной строке или терминалу

# Установка
Клонируйте репозиторий на ваше устройство:
```git
git clone https://github.com/GusevGrishaEm1/computer-club.git
cd ваш_репозиторий
```
# Запуск
```bash
go run main.go test.txt 
```

# Запуск через Docker
```bash
docker build -t computer_club .
docker run computer_club
```

# Пример использования сервиса
```go
package main

import (
	"computer_club/service"
	"fmt"
)

func ExampleNewComputerClub() {
	config := service.Config{
		FilePath: "test.txt",
	}
	club := service.NewComputerClub()
	club.SetInput(config)
	err := club.ProcessInput()
	if err != nil {
		panic(err)
	}
	for _, output := range club.GetOutput() {
		fmt.Println(output)
	}
	// Output:
	// 09:00
	// 08:48 1 client1
	// 08:48 13 NotOpenYet
	// 09:41 1 client1
	// 09:48 1 client2
	// 09:52 3 client1
	// 09:52 13 ICanWaitNoLonger
	// 09:54 2 client1 1
	// 10:25 2 client2 2
	// 10:58 1 client3
	// 10:59 2 client3 3
	// 11:30 1 client4
	// 11:35 2 client4 2
	// 11:35 13 PlaceIsBusy
	// 11:45 3 client4
	// 12:33 4 client1
	// 12:33 12 client4 1
	// 12:43 4 client2
	// 15:52 4 client4
	// 19:00 11 client3
	// 19:00
	// 1 70 05:58
	// 2 30 02:18
	// 3 90 08:01
}

```



