package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

func main() {
	// Configure serial connection settings
	config := &serial.Config{
		Name:        "/dev/ttyAMA0",
		Baud:        38400,
		ReadTimeout: time.Second,
	}

	// Open serial port
	port, err := serial.OpenPort(config)
	if err != nil {
		log.Fatalf("Failed to open port: %v", err)
	}
	defer port.Close()

	// Buffer for reading data
	readBuffer := make([]byte, 128)
	messageBuffer := []byte{}

	fmt.Println("Reading from serial port...")

	for {
		n, err := port.Read(readBuffer)
		if err != nil {
			// log.Printf("Error reading from serial port: %v", err)
			continue
		}

		if n > 0 {
			// Append new data to the message buffer
			messageBuffer = append(messageBuffer, readBuffer[:n]...)

			// Check for start and end delimiters (192) to identify a complete message
			startIndex := -1
			endIndex := -1

			for i, b := range messageBuffer {
				if b == 192 {
					if startIndex == -1 {
						startIndex = i
					} else {
						endIndex = i
						break
					}
				}
			}

			// If we found a complete message (from start to end delimiter)
			if startIndex != -1 && endIndex != -1 && endIndex > startIndex {
				// Extract the complete message
				payload := messageBuffer[startIndex+1 : endIndex]

				fmt.Print("Raw Data:\n")
				for i, b := range payload {
					fmt.Printf(" %01d: %#x - %08b\n", i, b, b)
				}

				buf := NewDataView(payload)
				frameType, _ := buf.GetUint8(0, true)
				commandID, _ := buf.GetUint8(1, true)
				sequence, _ := buf.GetUint8(2, true)

				status, _ := buf.GetUint8(5, true)
				param, _ := buf.GetUint8(7, true)

				log.Printf("FrameType: %#X, CommandID: %#X, Sequence: %#X, Status: %#X, Param %#X\n", frameType, commandID, sequence, status, param)

				// frame, err := ParseFrame(payload)
				// if err != nil {
				// 	fmt.Println("Error parsing frame:", err)
				// 	return
				// }

				// fmt.Printf("Parsed Frame:\nHeader: %+v\nPayload: %X\n", frame.Header, frame.Payload)

				// parseZigbeeMessage(payload)

				// Remove the processed message from the buffer
				messageBuffer = messageBuffer[endIndex+1:]
			}
		}
	}
}

type DataView struct {
	buffer *bytes.Reader
}

func NewDataView(data []byte) *DataView {
	return &DataView{buffer: bytes.NewReader(data)}
}

func (dv *DataView) GetUint16(offset int, littleEndian bool) (uint16, error) {
	dv.buffer.Seek(int64(offset), 0) // Seek to offset
	var value uint16
	var order binary.ByteOrder
	if littleEndian {
		order = binary.LittleEndian
	} else {
		order = binary.BigEndian
	}
	err := binary.Read(dv.buffer, order, &value)
	return value, err
}

func (dv *DataView) GetUint8(offset int, littleEndian bool) (uint8, error) {
	dv.buffer.Seek(int64(offset), 0) // Seek to offset
	var value uint8
	var order binary.ByteOrder
	if littleEndian {
		order = binary.LittleEndian
	} else {
		order = binary.BigEndian
	}
	err := binary.Read(dv.buffer, order, &value)
	return value, err
}

// FrameHeader represents the structure of a deCONZ frame header.
type FrameHeader struct {
	Length    uint8
	FrameType byte
	CommandID byte
	Sequence  byte
}

// Frame represents the full deCONZ frame.
type Frame struct {
	Header  FrameHeader
	Payload []byte
}

// ParseFrameHeader parses the deCONZ frame header.
func ParseFrameHeader(frame []byte) (FrameHeader, error) {
	if len(frame) < 3 {
		return FrameHeader{}, errors.New("frame too short for frame header")
	}

	reader := bytes.NewReader(frame)
	var header FrameHeader
	if err := binary.Read(reader, binary.LittleEndian, &header.Length); err != nil {
		return FrameHeader{}, fmt.Errorf("failed to read length: %w", err)
	}

	header.FrameType, _ = reader.ReadByte()
	header.CommandID, _ = reader.ReadByte()
	header.Sequence, _ = reader.ReadByte()

	fmt.Printf("Parsed Header: Length=%d, FrameType=%#X, CommandID=%#X, Sequence=%#X\n",
		header.Length, header.FrameType, header.CommandID, header.Sequence)

	return header, nil
}

// ParseFrame parses the entire deCONZ frame.
func ParseFrame(data []byte) (*Frame, error) {
	if len(data) < 3 {
		return nil, errors.New("data too short for frame")
	}

	header, err := ParseFrameHeader(data)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Data Length=%d, Expected Length=%d\n", len(data), header.Length)

	if len(data) < int(header.Length) {
		return nil, errors.New("data length mismatch")
	}

	payload := data[5:header.Length]
	return &Frame{
		Header:  header,
		Payload: payload,
	}, nil
}
