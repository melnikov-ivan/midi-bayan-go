package main

import (
	"time"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

// Буфер значения характеристики — глобальный, без аллокаций в прерывании (при чтении/записи BLE).
var charValueBuf [64]byte
var charValueLen int = 1 // начальное значение: 1 байт (0)
var hasNewValue bool     // флаг: в WriteEvent записали новое значение (вывод только по нему)

// StartBLEService включает адаптер, регистрирует сервис, запускает рекламу и блокируется.
// Вызывать из main в отдельной горутине: go StartBLEService().
func StartBLEService() {
	must(adapter.Enable())

	// 128-bit UUID'ы (произвольные)
	serviceUUID := bluetooth.NewUUID([16]byte{
		0x12, 0x34, 0x56, 0x78,
		0x12, 0x34,
		0x56, 0x78,
		0x12, 0x34,
		0x56, 0x78, 0x90, 0xab, 0xcd, 0xef,
	})

	charUUID := bluetooth.NewUUID([16]byte{
		0xfe, 0xdc, 0xba, 0x09,
		0x87, 0x65,
		0x43, 0x21,
		0x87, 0x65,
		0x43, 0x21, 0x10, 0x32, 0x54, 0x76,
	})

	// Переменная для хранения характеристики
	var char bluetooth.Characteristic

	// Начальное значение — срез глобального буфера, без аллокации при чтении клиентом
	charValueBuf[0] = 0

	// Регистрируем сервис через конфиг
	must(adapter.AddService(&bluetooth.Service{
		UUID: serviceUUID,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &char,
				UUID:   charUUID,
				Value:  charValueBuf[:charValueLen],
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicNotifyPermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					// В контексте прерывания нельзя делать аллокации (string, append, char.Write и т.д.).
					// Только копируем в глобальный буфер и ставим флаг.
					n := len(value)
					if n > len(charValueBuf) {
						n = len(charValueBuf)
					}
					for i := 0; i < n; i++ {
						charValueBuf[i] = value[i]
					}
					charValueLen = n
					hasNewValue = true
				},
			},
		},
	}))

	// Реклама: на nRF52 в TinyGo в объявлении поддерживаются только 16-битные UUID.
	// Берём короткую форму из первых двух байт 128-битного UUID сервиса (0x1234).
	advServiceUUID := bluetooth.New16BitUUID(0x1234)
	adv := adapter.DefaultAdvertisement()
	must(adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    "Midi-Bayan",
		ServiceUUIDs: []bluetooth.UUID{advServiceUUID},
	}))
	must(adv.Start())

	println("BLE peripheral started, advertising as 'Midi-Bayan'")

	// Бесконечный цикл: парсим входящие сообщения и синхронизируем только при новом значении
	for {
		time.Sleep(time.Second)
		if hasNewValue && charValueLen > 0 {
			hasNewValue = false
			msg := charValueBuf[:charValueLen]
			cmd, payload, ok := parseMessage(msg)
			if !ok {
				println("parse error or bad CRC, len=", charValueLen)
				continue
			}
			switch cmd {
			case cmdGetProgram:
				if ch, inst, vol, oct, ok := handleGetProgram(payload); ok {
					charValueBuf[0], charValueBuf[1], charValueBuf[2], charValueBuf[3] = ch, inst, vol, oct
					char.Write(charValueBuf[:4])
				}
			default:
				println("unknown command:", cmd)
			}
		}
		println("tick")
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
