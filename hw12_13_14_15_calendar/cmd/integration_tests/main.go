package main

import (
	"fmt"
	"os"

	"github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/integration"
)

func main() {
	if err := integration.RunIntegrationTests(); err != nil {
		fmt.Printf("Ошибка при выполнении интеграционных тестов: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
