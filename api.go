package main

const (
	cmdGetProgram byte = 0x01
)

const minMessageLen = 4 // cmd(1) + len(2) + crc(1), payload может быть 0

// crc8 считает CRC-8 (полином 0x07) по данным без последнего байта (место CRC).
func crc8(data []byte) byte {
	var crc byte = 0
	for _, b := range data {
		crc ^= b
		for i := 0; i < 8; i++ {
			if crc&0x80 != 0 {
				crc = (crc << 1) ^ 0x07
			} else {
				crc <<= 1
			}
		}
	}
	return crc
}

// parseMessage разбирает буфер: 1 байт команда, 2 байта длина payload (little-endian), payload, 1 байт CRC.
// Возвращает команду, payload и true при успехе; при ошибке ok == false.
func parseMessage(buf []byte) (cmd byte, payload []byte, ok bool) {
	if len(buf) < minMessageLen {
		return 0, nil, false
	}
	cmd = buf[0]
	payloadLen := int(buf[1]) | int(buf[2])<<8
	totalLen := 1 + 2 + payloadLen + 1
	if len(buf) < totalLen {
		return 0, nil, false
	}
	payload = buf[3 : 3+payloadLen]
	dataWithCrc := buf[:totalLen]
	gotCrc := dataWithCrc[totalLen-1]
	expectedCrc := crc8(dataWithCrc[:totalLen-1])
	if gotCrc != expectedCrc {
		return 0, nil, false
	}
	return cmd, payload, true
}

// handleGetProgram обрабатывает команду get_program: payload = [channel, instrument, octave].
// Instrument и octave для ответа берутся из config по указанному channel.
// Возвращает channel, instrument, octave и true при успехе; иначе ok == false.
func handleGetProgram(payload []byte) (channel, instrument, octave byte, ok bool) {
	if len(payload) != 3 {
		return 0, 0, 0, false
	}
	channel = payload[0]
	instrument, octave = GetChannelConfig(channel)
	println("get_program: channel=", channel, "instrument=", instrument, "octave=", octave)
	return channel, instrument, octave, true
}
