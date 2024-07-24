package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"unicode/utf16"
)

type TXT struct {
	XTranslation      uint16
	YTranslation      uint16
	Width             uint16
	Height            uint16
	TextIndex         uint16
	NumberOfTextLines uint16
	XScale            uint16
	YScale            uint16
	// Always seen to be 0, maybe not?
	_         uint32
	SomeValue uint32
	_         uint64
}

type TextLineTable struct {
	Width  uint16
	Height uint16
	_      uint16
	// TextSize is the size of the text in bytes
	TextSize uint16
}

func (r *Root) ParseTXT(data []byte) {
	reader := bytes.NewReader(data)
	var txt TXT
	err := binary.Read(reader, binary.LittleEndian, &txt)
	if err != nil {
		panic(err)
	}

	txtNode := TXTNode{
		XTranslation: txt.XTranslation,
		YTranslation: txt.YTranslation,
		Width:        txt.Width,
		Height:       txt.Height,
		TextIndex:    txt.TextIndex,
		XScale:       txt.XScale,
		YScale:       txt.YScale,
		SomeValue:    txt.SomeValue,
		Lines:        nil,
	}

	var lines []StringLines
	strOffset := 0
	for i := uint16(0); i < txt.NumberOfTextLines; i++ {
		var line TextLineTable
		err = binary.Read(reader, binary.LittleEndian, &line)
		if err != nil {
			panic(err)
		}

		str := utf16.Decode(r.txp[txt.TextIndex][strOffset : strOffset+(int(line.TextSize)/2)])
		strOffset += int(line.TextSize) / 2

		var builder strings.Builder
		startOfInvalid := true
		for _, _r := range str {
			// We have to filter out invalid characters. These may appear as placeholders for images.
			// The values are specific as well, therefore we must write them.
			if isValidXMLChar(_r) {
				if !startOfInvalid {
					builder.WriteString("]")
					startOfInvalid = true
				}

				builder.WriteRune(_r)
				continue
			}

			if startOfInvalid {
				builder.WriteString("[")
				startOfInvalid = false
			}

			builder.WriteString(fmt.Sprintf("%04x", binary.LittleEndian.Uint16(binary.LittleEndian.AppendUint16([]byte{}, uint16(_r)))))
		}

		lines = append(lines, StringLines{
			Width:  line.Width,
			Height: line.Height,
			String: builder.String(),
		})
	}

	txtNode.Lines = lines
	r.Panes = append(r.Panes, Children{TXT: &txtNode})
}

func (w *Writer) WriteTXT(txt *TXTNode) {
	header := SectionHeader{
		Type: SectionTypeTXT,
		Size: uint32(40 + (8 * len(txt.Lines))),
	}

	w.write(header)
	w.write(TXT{
		XTranslation:      txt.XTranslation,
		YTranslation:      txt.YTranslation,
		Width:             txt.Width,
		Height:            txt.Height,
		TextIndex:         txt.TextIndex,
		NumberOfTextLines: uint16(len(txt.Lines)),
		XScale:            txt.XScale,
		YScale:            txt.YScale,
		SomeValue:         txt.SomeValue,
	})

	for _, line := range txt.Lines {

		w.write(TextLineTable{
			Width:    line.Width,
			Height:   line.Height,
			TextSize: uint16(len(utf16.Encode([]rune(DecodeString(line.String)))) * 2),
		})
	}
}

func isValidXMLChar(r rune) bool {
	// Valid XML characters:
	// U+0009, U+000A, U+000D
	// U+0020 to U+D7FF
	// U+E000 to U+FFFD
	// U+10000 to U+10FFFF
	if r == 0x09 || r == 0x0A || r == 0x0D ||
		(r >= 0x20 && r <= 0xD7FF) ||
		(r >= 0xE000 && r <= 0xFFFD) ||
		(r >= 0x10000 && r <= 0x10FFFF) {
		return true
	}
	return false
}
