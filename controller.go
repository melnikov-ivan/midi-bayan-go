package main

// Пример отправки MIDI-команд через USB.
// USB MIDI поддерживается на платах с нативным USB (например Raspberry Pi Pico):
//
//   tinygo flash -target=xiao-ble controller.go
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

// KeyEvent — изменение состояния клавиши: бит 0..7, нажата или отпущена.
type KeyEvent struct {
	Bit    uint8 // индекс бита 0..7
	Pressed bool // true = нажата, false = отпущена
}

var led = machine.LED



func main() {
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.Low()

	m := midi.Port()
	println("USB MIDI Controller запущен")

	const (
		cable    = 0
		ch       = 1   // MIDI канал 1..16
		velocity = 100
	)

	keyCh := make(chan KeyEvent, 8)
	go RunKeyboard(keyCh)

	// Ноты берём из конфига keymap.BitToNote
	for ev := range keyCh {
		if int(ev.Bit) >= len(BitToNote) {
			continue
		}
		note := midi.Note(BitToNote[ev.Bit])
		if ev.Pressed {
			_ = m.NoteOn(cable, ch, note, velocity)
			println("MIDI: Note On ", note)
		} else {
			_ = m.NoteOff(cable, ch, note, 0)
			println("MIDI: Note Off", note)
		}
		blink()
	}
}

func blink() {
	led.High()
	time.Sleep(50 * time.Millisecond)
	led.Low()
}
