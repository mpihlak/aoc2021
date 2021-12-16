package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"unicode"
)

type BitStream struct {
	bits      []byte // 4 bit words
	bitPos    int    // bit position in the array
	maxBitPos int
	eof       bool
}

func NewBitStream(bits []byte) BitStream {
	return BitStream{bits, 0, len(bits) * 4, false}
}

// Take n bits
func (b *BitStream) Take(numBits int) *BitStream {
	res := BitStream{b.bits, b.bitPos, b.bitPos + numBits, b.eof}
	b.bitPos += numBits
	return &res
}

func (b *BitStream) AtEof() bool {
	return b.bitPos >= b.maxBitPos
}

func (b *BitStream) ReadBits(n int) int {
	result := int(0)
	for i := 0; i < n && b.bitPos < b.maxBitPos; i++ {
		arrayPos := b.bitPos / 4
		if arrayPos >= len(b.bits) {
			b.eof = true
			break
		}
		bitPos := b.bitPos % 4
		bit := b.bits[arrayPos] >> (3 - bitPos) & 1
		result = result<<1 | int(bit)
		b.bitPos += 1
	}

	return result
}

func ReadValue(b *BitStream) int {
	val := int(0)

	for {
		lengthBit := b.ReadBits(1)
		v := b.ReadBits(4)
		val = val<<4 | v

		if lengthBit == 0 {
			break
		}
	}

	return val
}

func SumPacketVersions(b *BitStream, maxPackets int) int {
	sum := int(0)
	packetCount := 0

	for !b.AtEof() {
		if maxPackets >= 0 && packetCount >= maxPackets {
			break
		}

		packetVersion := b.ReadBits(3)
		packetType := b.ReadBits(3)

		if packetType == 4 {
			// Literal value
			ReadValue(b)
		} else {
			// Operator packet
			lengthTypeId := b.ReadBits(1)
			if lengthTypeId == 0 {
				lengthInBits := int(b.ReadBits(15))
				newBits := b.Take(lengthInBits)
				sum += SumPacketVersions(newBits, -1)
			} else {
				lengthInPackets := int(b.ReadBits(11))
				sum += SumPacketVersions(b, lengthInPackets)
			}
		}

		sum += packetVersion
		packetCount += 1
	}

	return sum
}

func PacketValue(b *BitStream, maxPackets int, op int) int {
	values := []int{}
	packetCount := 0

	for !b.AtEof() {
		if maxPackets >= 0 && packetCount >= maxPackets {
			break
		}

		b.ReadBits(3) // Discard the packet version
		packetType := b.ReadBits(3)

		// Value of the packet we're currently looking at
		var currentValue int

		if packetType == 4 {
			currentValue = ReadValue(b)
		} else {
			lengthTypeId := b.ReadBits(1)
			lengthInPackets := -1
			lengthInBits := 0
			bits := b

			if lengthTypeId == 0 {
				lengthInBits = int(b.ReadBits(15))
				bits = b.Take(lengthInBits)
			} else {
				lengthInPackets = int(b.ReadBits(11))
			}

			currentValue = PacketValue(bits, lengthInPackets, packetType)
		}

		values = append(values, currentValue)
		packetCount += 1
	}

	// This could happen if run into EOF. Maybe there's a better way to handle it though
	if len(values) == 0 && op == 0 {
		return 0
	}

	var result int

	switch op {
	case 0:
		// sum of subpackets
		result = values[0]
		for _, v := range values[1:] {
			result += v
		}
	case 1:
		// product of subpackets
		result = values[0]
		for _, v := range values[1:] {
			result *= v
		}
	case 2:
		// MIN of subpackets
		result = math.MaxUint32
		for _, v := range values {
			if v < result {
				result = v
			}
		}
	case 3:
		// MAX of subpackets
		result = 0
		for _, v := range values {
			if v > result {
				result = v
			}
		}
	case 5:
		// Greater than: 1 if 1st subpacket is greater than the second. 0 otherwise
		if values[0] > values[1] {
			result = 1
		} else {
			result = 0
		}
	case 6:
		// Less than: 1 if 1st subpacket is less than the second. 0 otherwise
		if values[0] < values[1] {
			result = 1
		} else {
			result = 0
		}
	case 7:
		// Equal to: 1 of the subpackets are equal, 0 otherwise.
		if values[0] == values[1] {
			result = 1
		} else {
			result = 0
		}
	case -1:
		// no op, return the first value
		result = values[0]
	default:
		panic("invalid packet type")
	}

	return result
}

func main() {
	fileName := "input.txt"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	buf, _ := os.ReadFile(fileName)
	input := strings.TrimSuffix(string(buf), "\n")
	bits := make([]byte, len(input))

	for pos, c := range input {
		if unicode.IsDigit(c) {
			bits[pos] = byte(c - '0')
		} else {
			bits[pos] = byte(c-'A') + 10
		}
	}

	bs := NewBitStream(bits)

	sum := SumPacketVersions(&bs, -1)
	fmt.Println("Sum of versions =", sum)

	bs = NewBitStream(bits)
	val := PacketValue(&bs, -1, -1)
	fmt.Println("Value of the packet =", val)
}
