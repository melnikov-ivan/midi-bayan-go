package main

// tinygo flash -target=waveshare-rp2040-zero
// tinygo flash -target=xiao-ble keyboard.go
// tinygo monitor
import (
	"machine"
	"time"
)

var (
	led = machine.LED // 16

	// Пины для подключения 74HC165N
	// SH/LD (Shift/Load) - управляющий пин для загрузки параллельных данных
	shiftLoadPin = machine.D0

	// CLK (Clock) - тактовый сигнал для сдвига данных
	clockPin = machine.D1

	// QH (Serial Output) - последовательный выход данных
	dataPin = machine.D2
)

// readShiftRegister читает 8 бит данных из регистра сдвига 74HC165N.
// 74HC165 сдвигает по фронту CLK (LOW→HIGH). После сдвига нужно дать время
// выходу QH установиться (tpd ~20–30 ns), иначе читается старый/неверный бит.
const shiftDelay = 1 * time.Microsecond 

func readShiftRegister() uint8 {
	var data uint8 = 0

	// Шаг 1: Загрузка параллельных данных (SH/LD = LOW)
	shiftLoadPin.Low()
	time.Sleep(shiftDelay)

	// Шаг 2: Режим сдвига (SH/LD = HIGH) — первый бит уже на QH
	shiftLoadPin.High()
	time.Sleep(shiftDelay)

	// Шаг 3: Читаем 8 бит. Порядок: прочитать QH → такт → (новый бит на QH).
	// Первый бит уже на QH после load; остальные — после каждого сдвига.
	for i := 0; i < 8; i++ {
		if i > 0 {
			// Сдвиг: фронт LOW→HIGH. После сдвига на QH появляется следующий бит.
			clockPin.Low()
			time.Sleep(shiftDelay)
			clockPin.High()
			time.Sleep(shiftDelay) // время установки QH после сдвига
			clockPin.Low()
			time.Sleep(shiftDelay)
		}
		// Читаем бит с QH
		if dataPin.Get() {
			data |= (1 << (7 - i))
		}
	}

	return data
}

func main() {
	// Настройка LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.High()

	// Настройка пинов для 74HC165N
	// SH/LD - выход
	shiftLoadPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	shiftLoadPin.High() // По умолчанию в режиме сдвига

	// CLK - выход
	clockPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	clockPin.Low()

	// QH - вход
	dataPin.Configure(machine.PinConfig{Mode: machine.PinInput})

	println("Инициализация регистра сдвига 74HC165N")
	println("SH/LD: D0, CLK: D1, QH: D2")
	println("Начало чтения данных...")

	var n = 0
	// Основной цикл чтения данных
	for {
		// Читаем данные из регистра сдвига
		data := readShiftRegister()

		// Выводим прочитанные данные
		println("Данные:", formatBinary(data))

		n = n + 1
		if (n == 10000) {
			led.Low()
			time.Sleep(250 * time.Millisecond)
			led.High()
			n = 0
		}
		// Задержка перед следующим чтением
		// time.Sleep(50 * time.Millisecond)
		time.Sleep(1)
	}
}

// formatBinary форматирует байт в бинарную строку для удобного отображения
func formatBinary(b uint8) string {
	var result string
	for i := 7; i >= 0; i-- {
		if b&(1<<i) != 0 {
			result += "1"
		} else {
			result += "0"
		}
	}
	return result
}
