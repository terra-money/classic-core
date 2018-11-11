/*******************************************************************************
*   (c) 2018 ZondaX GmbH
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
********************************************************************************/

package ledger_goclient

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ZondaX/hid-go"
	"math"
)

const (
	VendorLedger    = 0x2c97
	UsagePageLedger = 0xffa0
	//ProductNano     = 1
	Channel    = 0x0101
	PacketSize = 64
	CLA        = 0x55

	INSGetVersion         = 0
	INSPublicKeySECP256K1 = 1
	INSPublicKeyED25519   = 2

	// Sign sdk.Msg (transaction part of the message)
	INSSignSECP256K1 = 3
	INSSignED25519   = 4

	// Sign sdk.StdSignMsg (full message)
	INSSignSECP256K1_StdSignMsg = 5
	INSSignED25519_StdSignMsg   = 6

	INSHash                   = 100
	INSPublicKeySECP256K1Test = 101
	INSPublicKeyED25519Test   = 102
	INSSignSECP256K1Test      = 103
	INSSignED25519Test        = 104

	MessageChunkSize = 250
)

type VersionInfo struct {
	AppId uint8
	Major uint8
	Minor uint8
	Patch uint8
}

type Ledger struct {
	device  Device
	Logging bool
}

func NewLedger(dev Device) *Ledger {
	return &Ledger{
		device:  dev,
		Logging: false,
	}
}

func ListDevices() {
	devices, err := hid.Devices()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	if len(devices) == 0 {
		fmt.Printf("No devices")
	}

	for _, d := range devices {
		fmt.Printf("============ %s\n", d.Path)
		fmt.Printf("Manufacturer  : %s\n", d.Manufacturer)
		fmt.Printf("Product       : %s\n", d.Product)
		fmt.Printf("ProductID     : %x\n", d.ProductID)
		fmt.Printf("VendorID      : %x\n", d.VendorID)
		fmt.Printf("VersionNumber : %x\n", d.VersionNumber)
		fmt.Printf("UsagePage     : %x\n", d.UsagePage)
		fmt.Printf("Usage         : %x\n", d.Usage)
		fmt.Printf("\n")
	}
}

func FindLedger() (*Ledger, error) {
	devices, err := hid.Devices()
	if err != nil {
		return nil, err
	}
	for _, d := range devices {
		if d.VendorID == VendorLedger && d.UsagePage == UsagePageLedger {
			device, err := d.Open()
			if err != nil {
				return nil, err
			}
			return NewLedger(device), nil
		}
	}
	return nil, errors.New("no ledger connected")
}

// A Device provides access to a HID device.
type Device interface {
	// Close closes the device and associated resources.
	Close()

	// Write writes an output report to device. The first byte must be the
	// report number to write, zero if the device does not use numbered reports.
	Write([]byte) error

	// ReadCh returns a channel that will be sent input reports from the device.
	// If the device uses numbered reports, the first byte will be the report
	// number.
	ReadCh() <-chan []byte

	// ReadError returns the read error, if any after the channel returned from
	// ReadCh has been closed.
	ReadError() error
}

func (ledger *Ledger) Exchange(command []byte) ([]byte, error) {
	if ledger.Logging {
		fmt.Printf("[%3d]=> %x\n", len(command), command)
	}

	if len(command) < 5 {
		return nil, fmt.Errorf("APDU commands should not be smaller than 5")
	}

	if (byte)(len(command)-5) != command[4] {
		return nil, fmt.Errorf("APDU[data length] mismatch")
	}

	serializedCommand, err := WrapCommandAPDU(Channel, command, PacketSize, false)

	if err != nil {
		return nil, err
	}

	// Write all the packets
	err = ledger.device.Write(serializedCommand[:PacketSize])
	if err != nil {
		return nil, err
	}
	for len(serializedCommand) > PacketSize {
		serializedCommand = serializedCommand[PacketSize:]
		err = ledger.device.Write(serializedCommand[:PacketSize])
		if err != nil {
			return nil, err
		}
	}

	input := ledger.device.ReadCh()
	response, err := UnwrapResponseAPDU(Channel, input, PacketSize, false)

	if len(response) < 2 {
		return nil, fmt.Errorf("lost connection")
	}

	swOffset := len(response) - 2
	sw := codec.Uint16(response[swOffset:])

	if ledger.Logging {
		fmt.Printf("Response: [%3d]<= %s\n", len(response[:swOffset]), response[:swOffset])
	}
	// FIXME: Code and description don't match for 0x6982 and 0x6983 based on
	// apdu spec: https://www.eftlab.co.uk/index.php/site-map/knowledge-base/118-apdu-response-list
	if sw != 0x9000 {
		switch sw {
		case 0x6400:
			return nil, errors.New("[APDU_CODE_EXECUTION_ERROR] No information given (NV-Ram not changed).")
		case 0x6700:
			return nil, errors.New("[APDU_CODE_WRONG_LENGTH] Wrong length.")
		case 0x6982:
			return nil, errors.New("[APDU_CODE_EMPTY_BUFFER] Security condition not satisfied.")
		case 0x6983:
			return nil, errors.New("[APDU_CODE_OUTPUT_BUFFER_TOO_SMALL] Authentication method blocked.")
		case 0x6984:
			return nil, errors.New("[APDU_CODE_DATA_INVALID] Referenced data reversibly blocked (invalidated).")
		case 0x6985:
			return nil, errors.New("[APDU_CODE_CONDITIONS_NOT_SATISFIED] Conditions of use not satisfied.")
		case 0x6986:
			return nil, errors.New("[APDU_CODE_COMMAND_NOT_ALLOWED] Command not allowed (no current EF).")
		case 0x6A80:
			return nil, errors.New("[APDU_CODE_BAD_KEY_HANDLE] The parameters in the data field are incorrect.")
		case 0x6B00:
			return nil, errors.New("[APDU_CODE_INVALIDP1P2] Wrong parameter(s) P1-P2.")
		case 0x6D00:
			return nil, errors.New("[APDU_CODE_INS_NOT_SUPPORTED] Instruction code not supported or invalid.")
		case 0x6E00:
			return nil, errors.New("[APDU_CODE_CLA_NOT_SUPPORTED] Class not supported.")
		case 0x6F00:
			return nil, errors.New("APDU_CODE_UNKNOWN")
		case 0x6F01:
			return nil, errors.New("APDU_CODE_SIGN_VERIFY_ERROR")
		}
		return nil, fmt.Errorf("invalid status %04x", sw)
	}

	return response[:swOffset], nil
}

func (ledger *Ledger) GetVersion() (*VersionInfo, error) {
	message := []byte{CLA, INSGetVersion, 0, 0, 0}
	response, err := ledger.Exchange(message)

	if err != nil {
		return nil, err
	}

	if len(response) < 4 {
		return nil, fmt.Errorf("invalid response")
	}

	return &VersionInfo{
		AppId: response[0],
		Major: response[1],
		Minor: response[2],
		Patch: response[3],
	}, nil
}

func getBip32bytes(bip32_path []uint32) ([]byte, error) {
	message := make([]byte, 41)
	if len(bip32_path) > 10 {
		return nil, fmt.Errorf("maximum bip32 depth = 10")
	}
	message[0] = byte(len(bip32_path))
	for index, element := range bip32_path {
		pos := 1 + index*4
		value := element
		// Harden 0, 1, 2
		if index <= 2 {
			value = 0x80000000 | element
		}
		binary.LittleEndian.PutUint32(message[pos:], value)
	}
	return message, nil
}

func (ledger *Ledger) sign(instruction byte, bip32_path []uint32, transaction []byte) ([]byte, error) {
	var packetIndex byte = 1
	var packetCount byte = 1 + byte(math.Ceil(float64(len(transaction))/float64(MessageChunkSize)))

	var finalResponse []byte

	var message []byte

	for packetIndex <= packetCount {
		chunk := MessageChunkSize
		if packetIndex == 1 {
			pathBytes, err := getBip32bytes(bip32_path)
			if err != nil {
				return nil, err
			}
			header := []byte{CLA, instruction, packetIndex, packetCount, byte(len(pathBytes))}
			message = append(header, pathBytes...)
		} else {
			if len(transaction) < MessageChunkSize {
				chunk = len(transaction)
			}
			header := []byte{CLA, instruction, packetIndex, packetCount, byte(chunk)}
			message = append(header, transaction[:chunk]...)
		}

		response, err := ledger.Exchange(message)
		if err != nil {
			return nil, err
		}

		finalResponse = response
		if packetIndex > 1 {
			transaction = transaction[chunk:]
		}
		packetIndex++

	}
	return finalResponse, nil
}

func (ledger *Ledger) SignSECP256K1(bip32_path []uint32, transaction []byte) ([]byte, error) {
	return ledger.sign(INSSignSECP256K1, bip32_path, transaction)
}

func (ledger *Ledger) SignED25519(bip32_path []uint32, transaction []byte) ([]byte, error) {
	return ledger.sign(INSSignED25519, bip32_path, transaction)
}

func (ledger *Ledger) SignSECP256K1_StdSignMsg(bip32_path []uint32, transaction []byte) ([]byte, error) {
	return ledger.sign(INSSignSECP256K1_StdSignMsg, bip32_path, transaction)
}

func (ledger *Ledger) SignED25519_StdSignMsg(bip32_path []uint32, transaction []byte) ([]byte, error) {
	return ledger.sign(INSSignED25519_StdSignMsg, bip32_path, transaction)
}

func (ledger *Ledger) GetPublicKeySECP256K1(bip32_path []uint32) ([]byte, error) {
	pathBytes, err := getBip32bytes(bip32_path)
	if err != nil {
		return nil, err
	}
	header := []byte{CLA, INSPublicKeySECP256K1, 0, 0, byte(len(pathBytes))}
	message := append(header, pathBytes...)

	response, err := ledger.Exchange(message)

	if err != nil {
		return nil, err
	}

	if len(response) < 4 {
		return nil, fmt.Errorf("invalid response")
	}

	return response, nil
}

func (ledger *Ledger) GetPublicKeyED25519(bip32_path []uint32) ([]byte, error) {
	pathBytes, err := getBip32bytes(bip32_path)
	if err != nil {
		return nil, err
	}

	header := []byte{CLA, INSPublicKeyED25519, 0, 0, byte(len(pathBytes))}
	message := append(header, pathBytes...)

	response, err := ledger.Exchange(message)

	if err != nil {
		return nil, err
	}

	if len(response) < 4 {
		return nil, fmt.Errorf("invalid response")
	}

	return response, nil
}

func (ledger *Ledger) Hash(transaction []byte) ([]byte, error) {

	var packetIndex = byte(1)
	var packetCount = byte(math.Ceil(float64(len(transaction)) / float64(MessageChunkSize)))

	var finalResponse []byte
	for packetIndex <= packetCount {
		chunk := MessageChunkSize
		if len(transaction) < MessageChunkSize {
			chunk = len(transaction)
		}

		header := []byte{CLA, INSHash, packetIndex, packetCount, byte(chunk)}
		message := append(header, transaction[:chunk]...)
		response, err := ledger.Exchange(message)

		if err != nil {
			return nil, err
		}
		finalResponse = response
		packetIndex++
		transaction = transaction[chunk:]
	}
	return finalResponse, nil
}

func (ledger *Ledger) TestGetPublicKeySECP256K1() ([]byte, error) {
	message := []byte{CLA, INSPublicKeySECP256K1Test, 0, 0, 0}
	response, err := ledger.Exchange(message)

	if err != nil {
		return nil, err
	}

	if len(response) < 4 {
		return nil, fmt.Errorf("invalid response")
	}

	return response, nil
}

func (ledger *Ledger) TestGetPublicKeyED25519() ([]byte, error) {
	message := []byte{CLA, INSPublicKeyED25519Test, 0, 0, 0}
	response, err := ledger.Exchange(message)

	if err != nil {
		return nil, err
	}

	if len(response) < 4 {
		return nil, fmt.Errorf("invalid response")
	}

	return response, nil
}

func (ledger *Ledger) TestSignSECP256K1(transaction []byte) ([]byte, error) {
	var packetIndex byte = 1
	var packetCount byte = byte(math.Ceil(float64(len(transaction)) / float64(MessageChunkSize)))

	var finalResponse []byte

	for packetIndex <= packetCount {

		chunk := MessageChunkSize
		if len(transaction) < MessageChunkSize {
			chunk = len(transaction)
		}

		header := []byte{CLA, INSSignSECP256K1Test, packetIndex, packetCount, byte(chunk)}
		message := append(header, transaction[:chunk]...)

		response, err := ledger.Exchange(message)

		if err != nil {
			return nil, err
		}

		finalResponse = response
		packetIndex++
		transaction = transaction[chunk:]
	}
	return finalResponse, nil
}

func (ledger *Ledger) TestSignED25519(transaction []byte) ([]byte, error) {
	var packetIndex byte = 1
	var packetCount byte = byte(math.Ceil(float64(len(transaction)) / float64(MessageChunkSize)))

	var finalResponse []byte

	for packetIndex <= packetCount {
		chunk := MessageChunkSize
		if len(transaction) < MessageChunkSize {
			chunk = len(transaction)
		}
		header := []byte{CLA, INSSignED25519Test, packetIndex, packetCount, byte(chunk)}
		message := append(header, transaction[:chunk]...)

		response, err := ledger.Exchange(message)

		if err != nil {
			return nil, err
		}

		finalResponse = response
		packetIndex++
		transaction = transaction[chunk:]
	}
	return finalResponse, nil
}
