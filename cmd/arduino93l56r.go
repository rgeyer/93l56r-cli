package cmd

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/dim13/cobs"
	"github.com/jacobsa/go-serial/serial"
)

type Arduino93L56R struct {
	serialOpts serial.OpenOptions
	serial     io.ReadWriteCloser
	reader     *bufio.Reader
}

func NewArduino93L56R(serPort string) *Arduino93L56R {
	return &Arduino93L56R{
		serialOpts: serial.OpenOptions{
			PortName:              serPort,
			BaudRate:              9600,
			DataBits:              8,
			StopBits:              1,
			InterCharacterTimeout: 200,
			MinimumReadSize:       0,
			ParityMode:            serial.PARITY_NONE,
		},
	}
}

func (a *Arduino93L56R) Connect() error {
	ser, err := serial.Open(a.serialOpts)
	if err != nil {
		return fmt.Errorf("Unable to open Serial Port: %s", err)
	}
	a.serial = ser
	a.reader = bufio.NewReader(a.serial)

	// c1 := make(chan error)
	// go func() {
	for i := 1; i <= 50; i++ {
		a.serial.Write(cobs.Encode([]byte{0x00}))

		readBytes, err := a.reader.ReadBytes(0x00)
		if err != nil && err == io.EOF {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if err != nil {
			return fmt.Errorf("Arduino did not respond as expected to reset request. Error: %s", err)
		}

		response := cobs.Decode(readBytes)
		if response[0] != 128 {
			return fmt.Errorf("Arduino acknowledged reset request with unexpected packet. Expected command 128, got %d", response[0])
		}
		break
	}

	return nil
}

func (a *Arduino93L56R) Close() {
	a.serial.Close()
}

// TODO: This only works with I2C EEPROMs with up to 256 addresses, since the
// arduino only sends a single address byte rather than two for a 16bit int
func (a *Arduino93L56R) I2CRead(addr int, length int) ([]byte, error) {
	addrMsb := byte(addr >> 8)
	addrLsb := byte(addr & 0xFF)

	lenMsb := byte(length >> 8)
	lenLsb := byte(length & 0xFF)

	readBuf := make([]byte, length)
	_, err := a.serial.Write(cobs.Encode([]byte{0x03, 0x50, addrMsb, addrLsb, lenMsb, lenLsb}))
	if err != nil {
		return readBuf, fmt.Errorf("Unable to send read request. Error: %s", err)
	}

	byteCnt, err := io.ReadAtLeast(a.serial, readBuf, length)
	if byteCnt != length || err != nil {
		return readBuf, fmt.Errorf("Did not receive all bytes from EEPROM read. Expected %d bytes, got %d. Error: %s\n\nResponse:\n%s", length, byteCnt, err, hex.Dump(readBuf[:byteCnt]))
	}

	return readBuf, nil
}

func (a *Arduino93L56R) Read(addr int, length int, icType string) ([]byte, error) {
	addrMsb := byte(addr >> 8)
	addrLsb := byte(addr & 0xFF)

	lenMsb := byte(length >> 8)
	lenLsb := byte(length & 0xFF)

	var rawBytes []byte
	if icType == "microwire" {
		rawBytes = []byte{0x01, addrMsb, addrLsb, lenMsb, lenLsb}
	}
	if icType == "i2c" {
		rawBytes = []byte{0x03, 0x50, addrMsb, addrLsb, lenMsb, lenLsb}
	}
	packetBytes := cobs.Encode(rawBytes)

	readBuf := make([]byte, length)
	_, err := a.serial.Write(packetBytes)
	if err != nil {
		return readBuf, fmt.Errorf("Unable to send read request. Error: %s", err)
	}

	byteCnt, err := io.ReadAtLeast(a.serial, readBuf, length)
	if byteCnt != length || err != nil {
		return readBuf, fmt.Errorf("Did not receive all bytes from EEPROM read. Expected %d bytes, got %d. Error: %s\n\nResponse:\n%s", length, byteCnt, err, hex.Dump(readBuf[:byteCnt]))
	}

	return readBuf, nil
}

func (a *Arduino93L56R) Write(addr int, buf []byte, icType string) error {
	var ackCmd byte
	addrMsb := byte(addr >> 8)
	addrLsb := byte(addr & 0xFF)

	lenMsb := byte(len(buf) >> 8)
	lenLsb := byte(len(buf) & 0xFF)
	// TODO: Length here is *actual* length in bytes. The eeprom has 16bit registers
	// so length is actually half of length of the supplied buffer. Everything
	// downstream does the work to translate it. Not sure if this should be register
	// length, rather than *actual* length.
	var rawBytes []byte
	if icType == "microwire" {
		rawBytes = append([]byte{0x02, addrMsb, addrLsb, lenMsb, lenLsb}, buf...)
		ackCmd = 130
	}
	if icType == "i2c" {
		rawBytes = append([]byte{0x04, 0x50, addrMsb, addrLsb, lenMsb, lenLsb}, buf...)
		ackCmd = 132
	}
	packetBytes := cobs.Encode(rawBytes)

	if len(packetBytes) > 64 {
		return fmt.Errorf("The resulting COBS packet for the write request exceeds 64 bytes and will overflow the Arduino Serial buffer. Actual size was %d", len(packetBytes))
	}

	wroteBytes, err := a.serial.Write(packetBytes)
	if wroteBytes != len(packetBytes) || err != nil {
		return fmt.Errorf("Unable to send buffer load request. Expected %d bytes written, got %d. Error: %s", len(packetBytes), wroteBytes, err)
	}

	writeResponseReceived := false
	for i := 1; i <= 50; i++ {
		readBytes, err := a.reader.ReadBytes(0x00)
		if err != nil && err == io.EOF {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		writeResponseReceived = true

		if err != nil {
			return fmt.Errorf("Arduino did not respond as expected to write request. Error: %s", err)
		}

		response := cobs.Decode(readBytes)
		if response[0] != ackCmd {
			return fmt.Errorf("Arduino acknowledged write request with unexpected packet. Expected command %d, got %d", ackCmd, response[0])
		}
		break
	}
	if writeResponseReceived {
		return nil
	} else {
		return fmt.Errorf("Timed out waiting for response to write request")
	}
}
