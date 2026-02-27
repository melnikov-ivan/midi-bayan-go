const CMD_GET_PROGRAM = 0x01;

function crc8(data) {
    let crc = 0;
    for (let i = 0; i < data.length; i++) {
        crc ^= data[i];
        for (let b = 0; b < 8; b++) {
            if (crc & 0x80) {
                crc = ((crc << 1) ^ 0x07) & 0xff;
            } else {
                crc = (crc << 1) & 0xff;
            }
        }
    }
    return crc;
}

function buildGetProgramMessage(channel) {
    const payload = new Uint8Array([channel, 0, 0]);
    const payloadLen = payload.length;
    const msg = new Uint8Array(1 + 2 + payloadLen + 1);
    msg[0] = CMD_GET_PROGRAM;
    msg[1] = payloadLen & 0xff;
    msg[2] = (payloadLen >> 8) & 0xff;
    msg.set(payload, 3);
    msg[3 + payloadLen] = crc8(msg.subarray(0, 3 + payloadLen));
    return msg;
}

function fillProgramFields(instrument, volume, octave) {
    const instEl = document.getElementById('instrumentValue');
    instEl.value = String(Number(instrument) & 0x7f);
    const volEl = document.getElementById('volumeValue');
    volEl.value = volume;
    document.getElementById('volumeValueDisplay').textContent = volume;
    const octaveClamped = Math.max(-3, Math.min(3, Number(octave)));
    document.getElementById('octaveValue').value = octaveClamped;
    document.getElementById('octaveValueDisplay').textContent = octaveClamped;
}
