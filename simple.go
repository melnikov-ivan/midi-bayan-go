package main

// BLE peripheral (XIAO BLE):
//   tinygo flash -target=xiao-ble -tags=ble simple.go
//   tinygo monitor
import (
	"tinygo.org/x/bluetooth"
	"time"
)

var adapter = bluetooth.DefaultAdapter

func main() {
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

	// Регистрируем сервис через конфиг
	must(adapter.AddService(&bluetooth.Service{
		UUID: serviceUUID,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &char,
				UUID:   charUUID,
				Value:  []byte{0}, // Начальное значение
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicNotifyPermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					// Обработка записи
					println("Received data:", string(value))
					// Обновляем значение характеристики
					char.Write(value)
				},
			},
		},
	}))

	// Реклама
	adv := adapter.DefaultAdvertisement()
	must(adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "XIAO-TinyGo",
	}))
	must(adv.Start())

	println("BLE peripheral started, advertising as 'XIAO-TinyGo'")

	// Бесконечный цикл для поддержания работы
	for {
		time.Sleep(time.Second)
		println("tick")
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
