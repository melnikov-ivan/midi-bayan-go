// out.go — отправка MIDI через UART (DIN, 31250 бод).
//
// MIDI по DIN использует UART 31250 бод. На XIAO BLE:
//   TX = D6 (P1_11), RX = D7 (P1_12)
//
// Прошивка:
//   tinygo flash -target=xiao-ble .
package main

import (
	"machine"
)

const (
	midiBaud = 31250
)

// startMidiOut инициализирует UART для MIDI. Вызывается из controller перед отправкой нот.
func startMidiOut() {
	machine.UART0.Configure(machine.UARTConfig{
		BaudRate: midiBaud,
		TX:       machine.UART_TX_PIN, // D6, P1_11
		RX:       machine.UART_RX_PIN, // D7, P1_12
	})
}

// SendNoteOn отправляет MIDI Note On по UART (0x90 | channel, note, velocity).
func SendNoteOn(channel uint8, note, velocity uint8) {
	ch := channel & 0x0F
	machine.UART0.Write([]byte{0x90 | ch, note & 0x7F, velocity & 0x7F})
}

// SendNoteOff отправляет MIDI Note Off по UART (0x80 | channel, note, 0).
func SendNoteOff(channel uint8, note uint8) {
	ch := channel & 0x0F
	machine.UART0.Write([]byte{0x80 | ch, note & 0x7F, 0})
}

// SendProgramChange отправляет MIDI Program Change по UART (0xC0 | channel, program).
func SendProgramChange(channel uint8, program uint8) {
	ch := channel & 0x0F
	machine.UART0.Write([]byte{0xC0 | ch, program & 0x7F})
}

// SendVolume отправляет MIDI Control Change #7 (Channel Volume) по UART (0xB0 | channel, 0x07, value).
func SendVolume(channel uint8, volume uint8) {
	ch := channel & 0x0F
	machine.UART0.Write([]byte{0xB0 | ch, 0x07, volume & 0x7F})
}
