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
	"github.com/pkg/errors"
	"encoding/binary"
)

var codec = binary.BigEndian

func SerializePacket(
	channel uint16,
	command []byte,
	packetSize int,
	sequenceIdx uint16,
	ble bool)	(result []byte, offset int, err error) {

	if packetSize < 3 {
		return nil, 0, errors.New("Packet size must be at least 3")
	}

	var headerOffset uint8

	result = make([]byte, packetSize)
	var buffer = result

	// Insert channel (2 bytes)
	if !ble {
		codec.PutUint16(buffer, channel)
		headerOffset += 2
	}

	// Insert tag (1 byte)
	buffer[headerOffset] = 0x05
	headerOffset += 1

	var commandLength uint16
	commandLength = uint16(len(command))

	// Insert sequenceIdx (2 bytes)
	codec.PutUint16(buffer[headerOffset:], sequenceIdx)
	headerOffset += 2

	// Only insert total size of the command in the first package
	if sequenceIdx == 0 {
		// Insert sequenceIdx (2 bytes)
		codec.PutUint16(buffer[headerOffset:], commandLength)
		headerOffset += 2
	}

	buffer = buffer[headerOffset:]
	offset = copy(buffer, command)
	return result, offset, nil
}

func DeserializePacket(
	channel uint16,
	buffer []byte,
	sequenceIdx uint16,
	ble bool)	(result []byte, totalResponseLength uint16, err error) {

	if (sequenceIdx == 0 && len(buffer) < 7) || (sequenceIdx > 0 && len(buffer) < 5) {
		return nil, 0, errors.New("Cannot deserialize the packet. Header information is missing.")
	}

	var headerOffset uint8

	if !ble {
		if codec.Uint16(buffer) != channel {
			return nil, 0, errors.New("Invalid channel")
		}
		headerOffset += 2
	}
	if buffer[headerOffset] != 0x05 {
		return nil, 0, errors.New("Invalid tag")
	}
	headerOffset++

	if codec.Uint16(buffer[headerOffset:]) != sequenceIdx {
		return nil, 0, errors.New("Wrong sequenceIdx")
	}
	headerOffset += 2

	if sequenceIdx == 0 {
		totalResponseLength = codec.Uint16(buffer[headerOffset:])
		headerOffset += 2
	}

	result = make([]byte, len(buffer) - int(headerOffset))
	copy(result, buffer[headerOffset:])

	return result, totalResponseLength, nil
}

// WrapCommandAPDU turns the command into a sequence of 64 byte packets
func WrapCommandAPDU(
	channel uint16,
	command []byte,
	packetSize int,
	ble bool) (result []byte, err error) {

	var offset int
	var totalResult []byte
	var sequenceIdx uint16
	for len(command) > 0 {
		result, offset, err = SerializePacket(channel, command, packetSize, sequenceIdx, ble)
		if err != nil {
			return nil, err
		}
		command = command[offset:]
		totalResult = append(totalResult, result...)
		sequenceIdx++
	}
	return totalResult, nil
}

// UnwrapResponseAPDU parses a response of 64 byte packets into the real data
func UnwrapResponseAPDU(channel uint16, pipe <- chan []byte, packetSize int, ble bool) ([]byte, error) {
	var sequenceIdx uint16

	var totalResult []byte
	var totalSize uint16
	var finished bool = false
	for !finished {

		// Read next packet from the channel
		buffer := <- pipe
		result, responseSize, err := DeserializePacket(channel, buffer, sequenceIdx, ble)
		if err != nil {
			return nil, err
		}
		if sequenceIdx == 0 {
			totalSize = responseSize
		}

		buffer = buffer[packetSize:]
		totalResult = append(totalResult, result...)
		sequenceIdx++

		if len(totalResult) >= int(totalSize) {
			finished = true
		}
	}

	// Remove trailing zeros
	totalResult = totalResult[:totalSize]
	return totalResult, nil
}