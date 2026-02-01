package main

// tinygo flash -target=xiao-ble usb.go
// tinygo monitor
import (
	// "io"
	"machine"
	"time"
)

var (
	led  = machine.LED
	uart = machine.UART0
)

func main() {
	// Настройка LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.Low()

	// Настройка последовательного порта (UART)
	uart.Configure(machine.UARTConfig{
		BaudRate: 9600,
		TX:       machine.UART_TX_PIN,
		RX:       machine.UART_RX_PIN,
	})

	println("Последовательный порт инициализирован")
	println("Отправьте '1' для включения LED или '0' для выключения")

	// Буфер для чтения одного байта
	buf := make([]byte, 1)

	// Бесконечный цикл чтения из последовательного порта
	for {
		// Пытаемся прочитать один байт
		// Read будет блокироваться до получения данных или вернет ошибку
		n, err := uart.Read(buf)
		if err != nil {
			// Если ошибка (например, нет данных), ждем немного и продолжаем
			time.Sleep(10 * time.Millisecond)
			continue
		}

		if n > 0 {
			char := buf[0]

			// Обработка символов '1' и '0'
			if char == '1' {
				led.High()
				println("LED включен")
			} else if char == '0' {
				led.Low()
				println("LED выключен")
			} else if char == '\n' || char == '\r' {
				// Игнорируем символы новой строки
				continue
			} else {
				// Неизвестный символ
				println("Получен неизвестный символ:", char)
			}
		}
	}
}
