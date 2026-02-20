package main

// Config хранит настройки канала: номер канала, инструмент и октаву.
type Config struct {
	Channel    byte
	Instrument byte
	Octave     byte
}

// Конфиги по каналам (0–15). Индекс = номер канала.
var ChannelConfigs [16]Config

func init() {
	for i := 0; i < 16; i++ {
		ChannelConfigs[i] = Config{
			Channel:    byte(i),
			Instrument: byte(i), // Acoustic Grand Piano
			Octave:     4, // средняя октава
		}
	}
}

// GetChannelConfig возвращает Instrument и Octave для канала channel.
// Если channel >= 16, возвращает 0, 0.
func GetChannelConfig(channel byte) (instrument, octave byte) {
	if channel >= 16 {
		return 0, 0
	}
	c := ChannelConfigs[channel]
	return c.Instrument, c.Octave
}
