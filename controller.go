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

// Пины для 74HC165N
var (
	shiftLoadPin = machine.D0
	clockPin     = machine.D1
	dataPin      = machine.D2
)

const shiftDelay = 1 * time.Microsecond

func readShiftRegister() uint8 {
	var data uint8
	shiftLoadPin.Low()
	time.Sleep(shiftDelay)
	shiftLoadPin.High()
	time.Sleep(shiftDelay)
	for i := 0; i < 8; i++ {
		if i > 0 {
			clockPin.Low()
			time.Sleep(shiftDelay)
			clockPin.High()
			time.Sleep(shiftDelay)
			clockPin.Low()
			time.Sleep(shiftDelay)
		}
		if dataPin.Get() {
			data |= (1 << (7 - i))
		}
	}
	return data
}

// RunKeyboard читает регистр сдвига, при изменении бита отправляет KeyEvent в канал ch.
func RunKeyboard(ch chan<- KeyEvent) {
	shiftLoadPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	shiftLoadPin.High()
	clockPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	clockPin.Low()
	dataPin.Configure(machine.PinConfig{Mode: machine.PinInput})
	var prev uint8
	for {
		data := readShiftRegister()
		for i := uint8(0); i < 8; i++ {
			mask := uint8(1 << i)
			was := (prev & mask) != 0
			now := (data & mask) != 0
			if was != now {
				ch <- KeyEvent{Bit: i, Pressed: now}
			}
		}
		prev = data
		time.Sleep(1)
	}
}

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

	// Читаем события из канала и отправляем MIDI: бит 0 = C4, бит 1 = C#4, ...
	for ev := range keyCh {
		note := midi.Note(uint8(midi.C4) + ev.Bit)
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
