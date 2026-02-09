# MIDI-клавиатура (controller.go уже включает код клавиатуры 74HC165)
TARGET ?= xiao-ble

flash:
	tinygo flash -target=$(TARGET) .

monitor:
	tinygo monitor

.PHONY: flash monitor
