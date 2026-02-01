package main

import (
	"machine"
	"time"

	"tinygo.org/x/bluetooth"
)

var (
	adapter = bluetooth.DefaultAdapter
	led     = machine.LED
)

// 12345678-9abc-def0-1234-56789abcdef0
var serviceUUID = bluetooth.NewUUID([16]byte{
    0xf0, 0xde, 0xbc, 0x9a,
    0x78, 0x56,
    0x34, 0x12,
	0xf0, 0xde,
    0xbc, 0x9a,
    0x78, 0x56, 0x34, 0x12,
})

// 12345678-9abc-def0-1234-56789abcdef1
var charUUID = bluetooth.NewUUID([16]byte{
    0xf1, 0xde, 0xbc, 0x9a,
    0x78, 0x56,
    0x34, 0x12,
	0xf0, 0xde,
    0xbc, 0x9a,
    0x78, 0x56, 0x34, 0x12,
})

func main() {
	// Настройка LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.Low()

	// Включить BLE адаптер
	must("enable BLE stack", adapter.Enable())

	// Настройка рекламы
	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    "XIAO-BLE Example",
		ServiceUUIDs: []bluetooth.UUID{serviceUUID},
	}))

	// Начать рекламу
	must("start adv", adv.Start())

	// Добавить сервис с характеристикой
	var char bluetooth.Characteristic
	must("add service", adapter.AddService(&bluetooth.Service{
		UUID: serviceUUID,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &char,
				UUID:   charUUID,
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicNotifyPermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					// Обработка записи - мигаем LED
					println("Received data:", string(value))
					led.High()
					time.Sleep(time.Millisecond * 200)
					led.Low()
				},
			},
		},
	}))

	// Устанавливаем начальное значение для чтения
	char.Write([]byte("Hello from XIAO-BLE!"))

	println("BLE Peripheral запущен. Имя устройства: XIAO-BLE Example")
	println("Ожидание подключения...")

	// Бесконечный цикл
	for {
		time.Sleep(1 * time.Second)
		println("heartbeat")
	}
}

func must(action string, err error) {
	if err != nil {
		println("failed to", action, ":", err.Error())
		for {
			time.Sleep(time.Minute)
		}
	}
}