package main

// MIDI-команды отправляются через UART (DIN, 31250 бод). См. out.go.
//
//   tinygo flash -target=xiao-ble .
//   tinygo monitor
import (
	"machine"
	"time"
)

// KeyEventType — тип события: нота (клавиша) или смена программы.
type KeyEventType uint8

const (
	NoteOn        KeyEventType = iota // событие клавиши (Bit, Pressed)
	ProgramChange                    // смена инструмента (Channel, Program)
)

// Event — изменение состояния клавиши (Type=NoteOn) либо событие Program Change (Type=ProgramChange).
type Event struct {
	Type   KeyEventType // NoteOn или ProgramChange
	Bit    uint8        // индекс бита 0..7 (для NoteOn)
	Pressed bool        // true = нажата, false = отпущена (для NoteOn)

	Channel uint8 // для ProgramChange
	Program uint8 // для ProgramChange
}

var led = machine.LED

// EventChannel — общий канал событий: клавиатура (RunKeyboard) и BLE (handleSetProgram Program Change).
var EventChannel chan Event

func main() {
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.Low()

	startMidiOut()
	println("MIDI Controller (UART) запущен")

	EventChannel = make(chan Event, 8)
	go StartBLEService()
	go RunKeyboard(EventChannel)

	const (
		ch       = 0   // MIDI канал 0..15
		velocity = 100
	)

	// Ноты берём из конфига keymap.BitToNote; Program Change приходит из BLE (handleSetProgram).
	for ev := range EventChannel {
		switch ev.Type {
		case ProgramChange:
			SendProgramChange(ev.Channel, ev.Program)
			println("MIDI: Program Change ch=", ev.Channel, "program=", ev.Program)
			blink()
		case NoteOn:
			if int(ev.Bit) >= len(BitToNote) {
				break
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
		default:
			// неизвестный тип — пропускаем
		}
	}
}

func blink() {
	led.High()
	time.Sleep(50 * time.Millisecond)
	led.Low()
}
