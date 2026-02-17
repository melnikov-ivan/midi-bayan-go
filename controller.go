package main

// MIDI-команды отправляются через UART (DIN, 31250 бод). См. out.go.
//
//   tinygo flash -target=xiao-ble .
//   tinygo monitor
import (
	"machine"
	"time"
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

	startMidiOut()
	println("MIDI Controller (UART) запущен")

	go StartBLEService()

	const (
		ch       = 0   // MIDI канал 0..15
		velocity = 100
	)

	keyCh := make(chan KeyEvent, 8)
	go RunKeyboard(keyCh)

	// Ноты берём из конфига keymap.BitToNote
	for ev := range keyCh {
		if int(ev.Bit) >= len(BitToNote) {
			continue
		}
		note := BitToNote[ev.Bit]
		if ev.Pressed {
			SendNoteOn(ch, note, velocity)
			println("MIDI: Note On ", note)
		} else {
			SendNoteOff(ch, note)
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
