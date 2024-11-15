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
	// messageBuffer := []byte{}

	fmt.Println("Reading from serial port...")
	for {
		n, err := port.Read(readBuffer)
		if err != nil {
			// log.Printf("Error reading from serial port: %v", err)
			continue
		}

		if n > 0 {
			frame, err := ParseFrame(readBuffer[:n])
			if err != nil {
				fmt.Println("Error parsing frame:", err)
				return
			}

			fmt.Printf("Parsed Frame:\nHeader: %+v\nPayload: %X\n", frame.Header, frame.Payload)
			// // Append new data to the message buffer
			// messageBuffer = append(messageBuffer, readBuffer[:n]...)

			// // Check for start and end delimiters (192) to identify a complete message
			// startIndex := -1
			// endIndex := -1

			// for i, b := range messageBuffer {
			// 	if b == 192 {
			// 		if startIndex == -1 {
			// 			startIndex = i
			// 		} else {
			// 			endIndex = i
			// 			break
			// 		}
			// 	}
			// }

			// // If we found a complete message (from start to end delimiter)
			// if startIndex != -1 && endIndex != -1 && endIndex > startIndex {
			// 	// Extract the complete message
			// 	payload := messageBuffer[startIndex : endIndex+1]

			// 	// fmt.Print("Raw Data:\n")
			// 	// for i, b := range payload {
			// 	// 	fmt.Printf(" %01d: %#x - %08b\n", i, b, b)
			// 	// }

			// 	// payload = []byte{0x0A, 0x00, 0x01, 0x02, 0x03, 0xAA, 0xBB, 0xCC}

			// 	frame, err := ParseFrame(payload)
			// 	if err != nil {
			// 		fmt.Println("Error parsing frame:", err)
			// 		return
			// 	}

			// 	fmt.Printf("Parsed Frame:\nHeader: %+v\nPayload: %X\n", frame.Header, frame.Payload)

			// 	// parseZigbeeMessage(payload)

			// 	// Remove the processed message from the buffer
			// 	messageBuffer = messageBuffer[endIndex+1:]
			// }
		}
	}
}

// Parsing function for Zigbee message structure
func parseZigbeeMessage(data []byte) {
	// Extract fields based on observed pattern
	startByte := data[0]
	endByte := data[len(data)-1]

	messageType := data[1]           // 2nd byte, possible message type or command
	deviceID := data[2]              // 3rd byte, could be device identifier
	payload := data[3 : len(data)-1] // Remaining data except start and end bytes

	// Interpret payload as fields or a single integer (depends on message format)
	fmt.Printf("Parsed Message:\n")
	fmt.Printf("  Start Byte: 0x%02X\n", startByte)
	fmt.Printf("  Message Type: 0x%02X\n", messageType)
	fmt.Printf("  Device ID: 0x%02X\n", deviceID)
	fmt.Print("  Payload:\n")
	for i, b := range payload {
		fmt.Printf("    Byte %d: 0x%02X", i, b)
	}
	fmt.Printf("  End Byte: 0x%02X\n\n", endByte)
}

// FrameHeader represents the structure of a deCONZ frame header.
type FrameHeader struct {
	Length    uint16
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
func ParseFrameHeader(data []byte) (FrameHeader, error) {
	if len(data) < 5 {
		return FrameHeader{}, errors.New("data too short for frame header")
	}

	reader := bytes.NewReader(data)
	var header FrameHeader
	if err := binary.Read(reader, binary.LittleEndian, &header.Length); err != nil {
		return FrameHeader{}, err
	}

	header.FrameType, _ = reader.ReadByte()
	header.CommandID, _ = reader.ReadByte()
	header.Sequence, _ = reader.ReadByte()

	return header, nil
}

// ParseFrame parses the entire deCONZ frame.
func ParseFrame(data []byte) (*Frame, error) {
	if len(data) < 5 {
		return nil, errors.New("data too short for frame")
	}

	header, err := ParseFrameHeader(data)
	if err != nil {
		return nil, err
	}

	if len(data) < int(header.Length) {
		return nil, errors.New("data length mismatch")
	}

	payload := data[5:header.Length]
	return &Frame{
		Header:  header,
		Payload: payload,
	}, nil
}
