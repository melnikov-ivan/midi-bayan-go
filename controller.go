package main

// Пример отправки MIDI-команд через USB.
// USB MIDI поддерживается на платах с нативным USB (например Raspberry Pi Pico):
//
//   tinygo flash -target=pico controller.go
//   tinygo monitor
//
// Подключите плату по USB к ПК — в системе появится USB MIDI-устройство.
// Для XIAO BLE (nRF52840) проверьте поддержку USB MIDI в TinyGo для вашей платы.
//
// Типичные команды: CC (Control Change), Note On, Note Off.
// Импорт: machine/usb/adc/midi (TinyGo 0.40+) или machine/usb/midi на других версиях.
import (
	"machine"
	"time"

	"machine/usb/adc/midi"
)

var led = machine.LED

func main() {
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.Low()

	// Создаём USB MIDI устройство (при первом вызове инициализирует USB стек)
	m := midi.Port()

	println("USB MIDI Controller запущен")
	println("Нажатие и отпускание ноты каждые 2 секунды...")

	const (
		cable    = 0
		channel  = 1   // MIDI канал 1..16
		velocity = 100
	)

	for {
		// Нажатие ноты (Note On)
		_ = m.NoteOn(cable, channel, midi.C4, velocity)
		println("MIDI: Note On  C4")
		blink()

		time.Sleep(500 * time.Millisecond)

		// Отпускание ноты (Note Off)
		_ = m.NoteOff(cable, channel, midi.C4, 0)
		println("MIDI: Note Off C4")
		blink()

		time.Sleep(1 * time.Second)

		// --- Примеры CC (закомментированы) ---
		// m.ControlChange(0, 1, 1, 64)   // CC#1 Modulation
		// m.ControlChange(0, 1, 7, 100)  // CC#7 Volume
		// m.ControlChange(0, 1, 10, 64)  // CC#10 Pan
	}
}

func blink() {
	led.High()
	time.Sleep(50 * time.Millisecond)
	led.Low()
}
