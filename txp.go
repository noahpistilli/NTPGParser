package main

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"
	"unicode/utf16"
)

func (r *Root) ParseTXP(data []byte) {
	reader := bytes.NewReader(data)
	var numberOfStrings uint32
	err := binary.Read(reader, binary.LittleEndian, &numberOfStrings)
	if err != nil {
		panic(err)
	}

	offsets := make([]uint32, numberOfStrings+1)
	for i := uint32(0); i < numberOfStrings+1; i++ {
		err = binary.Read(reader, binary.LittleEndian, &offsets[i])
		if err != nil {
			panic(err)
		}
	}

	// Read the strings now.
	strings := make([][]uint16, numberOfStrings)
	for i := uint32(0); i < numberOfStrings; i++ {
		size := offsets[i+1] - offsets[i]
		str := make([]byte, size)
		reader.Read(str)

		strings[i] = make([]uint16, size/2)
		err = binary.Read(bytes.NewReader(str), binary.LittleEndian, &strings[i])
		if err != nil {
			panic(err)
		}
	}

	// First string is a header and should be written to the TXP node.
	header := string(utf16.Decode(strings[0]))
	r.TXP = TXPNode{Content: header[:len(header)-1]}

	r.txp = strings
}

/*func ConvertUTF16ToString(data []uint16) string {

}*/

func (w *Writer) WriteTXP() {
	header := SectionHeader{
		Type: SectionTypeTXP,
		Size: uint32(12 + len(w.root.NAP.Contents)*4),
	}

	// We store TXP data in the text line nodes in the XML for easier legibility.
	var strs [][]uint16
	strs = append(strs, utf16.Encode([]rune(w.root.TXP.Content)))
	for _, pane := range w.root.Panes {
		if pane.TXT != nil {
			// Formulate the full string. Strings that are on different lines from each other are not seperated by
			// anything.
			str := ""
			for _, line := range pane.TXT.Lines {
				str += DecodeString(line.String)
			}

			// The len function behaves weirdly with Japanese characters as a string. Converting back to UTF-16
			strs = append(strs, utf16.Encode([]rune(str)))
		}
	}

	// Now formulate the node.
	offsets := make([]uint32, len(strs)+1)

	// First offset is relative to the number of offsets.
	size := 12 + uint32((len(strs)+1)*4)

	offsets[0] = uint32((len(strs) + 1) * 4)
	for i := 0; i < len(strs); i++ {
		// Add extra 2 bytes for null terminator (UTF-16)
		offsets[i+1] = offsets[i] + uint32(len(strs[i])*2) + 2
		size += uint32(len(strs[i])*2) + 2
	}

	for (w.Len()+int(size))%4 != 0 {
		size += 1
	}

	header.Size = size
	w.write(header)
	w.write(uint32(len(strs)))
	w.write(offsets)

	for _, str := range strs {
		// Decode the string as it may have the holder for invalid characters.
		w.write(str)
		w.write(uint16(0))
	}

	// Pad to 4 bytes
	for w.Len()%4 != 0 {
		w.WriteByte(0)
	}
}

func DecodeString(str string) string {
	var builder strings.Builder
	currInvalid := false
	for i := 0; i < len(str); i++ {
		if str[i] == '[' {
			currInvalid = true

			// We are guaranteed to have 4 bytes of data.
			strVal := string(str[i+1]) + string(str[i+2]) + string(str[i+3]) + string(str[i+4])
			value, _ := strconv.ParseUint(strVal, 16, 16)

			builder.WriteRune(rune(value))
			i += 4
			continue
		} else if str[i] == ']' {
			currInvalid = false
			continue
		}

		if currInvalid {
			strVal := string(str[i]) + string(str[i+1]) + string(str[i+2]) + string(str[i+3])
			value, _ := strconv.ParseUint(strVal, 16, 16)

			builder.WriteRune(rune(value))
			i += 3
		} else {
			builder.WriteByte(str[i])
		}
	}

	return builder.String()
}
