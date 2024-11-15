package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

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

// ExampleUsage demonstrates parsing a sample deCONZ frame.
func ExampleUsage() {
	// Example raw deCONZ frame (replace with actual data).
	rawData := []byte{0x0A, 0x00, 0x01, 0x02, 0x03, 0xAA, 0xBB, 0xCC}

	frame, err := ParseFrame(rawData)
	if err != nil {
		fmt.Println("Error parsing frame:", err)
		return
	}

	fmt.Printf("Parsed Frame:\nHeader: %+v\nPayload: %X\n", frame.Header, frame.Payload)
}
