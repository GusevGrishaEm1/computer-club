package main

import (
	"bufio"
	"computer_club/service"
	"fmt"
	"os"
)

// Компьютерный клуб
type ComputerClub interface {
	SetInput(c service.Config) error
	ProcessInput() error
	GetOutput() []string
}

// Точка входа в программу
func main() {
	config := service.NewConfig(os.Args[1])
	club := service.NewComputerClub()
	err := club.SetInput(config)
	if err != nil {
		panic(err)
	}
	err = club.ProcessInput()
	if err != nil {
		panic(err)
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for _, output := range club.GetOutput() {
		fmt.Fprintln(out, output)
	}
}
