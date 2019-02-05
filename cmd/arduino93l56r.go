package cmd

import (
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

type Arduino93L56R struct {
	serialOpts serial.OpenOptions
	serial     io.ReadWriteCloser
	eepromBuf  []byte
	cmdRespBuf []byte
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
		eepromBuf:  make([]byte, 256),
		cmdRespBuf: make([]byte, 4),
	}
}

func (a *Arduino93L56R) Connect() error {
	ser, err := serial.Open(a.serialOpts)
	if err != nil {
		return fmt.Errorf("Unable to open Serial Port: %s", err)
	}
	a.serial = ser

	a.serial.Write([]byte{'0', '\n'})
	// The sketch blinks the LED for 2 seconds
	time.Sleep(2 * time.Second)
	readBytes, err := io.ReadAtLeast(a.serial, a.eepromBuf, 5)
	resetBytes := make([]byte, 0)
	if readBytes <= 5 {
		resetBytes = a.eepromBuf[readBytes-5 : readBytes]
	}
	if err != nil || !reflect.DeepEqual(resetBytes, []byte{0x52, 0x65, 0x73, 0x65, 0x74}) {
		return fmt.Errorf("Arduino did not send expected Reset on initialization. Expected 5 bytes, got %d. Error: %s\n\nResponse: %s", readBytes, err, a.eepromBuf[:readBytes])
	}
	return nil
}

func (a *Arduino93L56R) Close() {
	a.serial.Close()
}

func (a *Arduino93L56R) Read() ([]byte, error) {
	wroteBytes, err := a.serial.Write([]byte{'r', '\n'})
	if wroteBytes != 2 || err != nil {
		return a.eepromBuf, fmt.Errorf("Unable to send read request. Expected 2 bytes written, got %d. Error: %s", wroteBytes, err)
	}

	readBytes, err := io.ReadFull(a.serial, a.cmdRespBuf)
	if readBytes != 4 || err != nil || !reflect.DeepEqual(a.cmdRespBuf, []byte{0x52, 0x45, 0x41, 0x44}) {
		return a.eepromBuf, fmt.Errorf("Did not receive expected response to read request. Expected 4 bytes, got %d. Error: %s\n\nResponse: %s", readBytes, err, a.cmdRespBuf[:readBytes])
	}

	readBytes, err = io.ReadAtLeast(a.serial, a.eepromBuf, 256)
	if readBytes != 256 || err != nil {
		return a.eepromBuf, fmt.Errorf("Did not receive all bytes from EEPROM read. Expected 256 bytes, got %d. Error: %s\n\nResponse: %s", readBytes, err, a.eepromBuf[:readBytes])
	}

	return a.eepromBuf, nil
}

func (a *Arduino93L56R) LoadBuffer(buf []byte) error {
	wroteBytes, err := a.serial.Write(append([]byte{'l', '\n'}, buf...))
	if wroteBytes != 2+256 || err != nil {
		return fmt.Errorf("Unable to send buffer load request. Expected 258 bytes written, got %d. Error: %s", wroteBytes, err)
	}

	readBytes, err := io.ReadFull(a.serial, a.cmdRespBuf)
	if readBytes != 4 || err != nil || !reflect.DeepEqual(a.cmdRespBuf, []byte{0x4c, 0x4f, 0x41, 0x44}) {
		return fmt.Errorf("Did not receive expected response to buffer load request. Expected 4 bytes, got %d. Error: %s\n\nResponse: %s", readBytes, err, a.cmdRespBuf[:readBytes])
	}
	return nil
}

func (a *Arduino93L56R) ReadBuffer() ([]byte, error) {
	wroteBytes, err := a.serial.Write([]byte{'p', '\n'})
	if wroteBytes != 2 || err != nil {
		return a.eepromBuf, fmt.Errorf("Unable to send buffer read request. Expected 2 bytes written, got %d. Error: %s", wroteBytes, err)
	}

	readBytes, err := io.ReadFull(a.serial, a.cmdRespBuf)
	if readBytes != 4 || err != nil || !reflect.DeepEqual(a.cmdRespBuf, []byte{0x50, 0x52, 0x4e, 0x54}) {
		return a.eepromBuf, fmt.Errorf("Did not receive expected response to buffer read request. Expected 4 bytes, got %d. Error: %s\n\nResponse: %s", readBytes, err, a.cmdRespBuf[:readBytes])
	}

	readBytes, err = io.ReadAtLeast(a.serial, a.eepromBuf, 256)
	if readBytes != 256 || err != nil {
		return a.eepromBuf, fmt.Errorf("Did not receive all bytes from buffer read. Expected 256 bytes, got %d. Error: %s\n\nResponse: %s", readBytes, err, a.eepromBuf[:readBytes])
	}

	return a.eepromBuf, nil
}

func (a *Arduino93L56R) WriteBuffer() error {
	wroteBytes, err := a.serial.Write([]byte{'w', '\n'})
	if wroteBytes != 2 || err != nil {
		return fmt.Errorf("Unable to send buffer write request. Expected 2 bytes written, got %d. Error: %s", wroteBytes, err)
	}

	// TODO: How to properly wait for writing to be completed?
	time.Sleep(4 * time.Second)

	readBytes, err := io.ReadFull(a.serial, a.cmdRespBuf)
	if readBytes != 4 || err != nil || !reflect.DeepEqual(a.cmdRespBuf, []byte{0x57, 0x52, 0x49, 0x54}) {
		return fmt.Errorf("Did not receive expected response to buffer write request. Expected 4 bytes, got %d. Error: %s\n\nResponse: %s", readBytes, err, a.cmdRespBuf[:readBytes])
	}

	return nil
}
