код написан на tinygo.

1. создаем модуль
`go mod init blinky`

2. пишем код 
в `main.go`

3. загружаем библиотеку
`go get tinygo.org/x/bluetooth`

4. прошиваем
переводим плату в boot mode двойным нажатием
`tinygo flash -target=xiao-ble -tags=ble main.go`

5. сборка
`tinygo build -target=xiao-ble -o firmware.uf2` посмотреть размер бинарника

6. отладка
читаем логи в терминале 
`tinygo monitor`