package main

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"log"
	"os"
)

// Header represents the header of our NTPG
type Header struct {
	Magic        [4]byte
	BOM          uint32
	FileSize     uint32
	HeaderLen    uint16
	SectionCount uint16
}

// SectionTypes are known parts of a NTPG.
type SectionTypes [4]byte

type SectionHeader struct {
	Type SectionTypes
	Size uint32
}

var (
	headerMagic    SectionTypes = [4]byte{'N', 'T', 'P', 'G'}
	SectionTypeNAP SectionTypes = [4]byte{'n', 'a', 'p', '1'}
	SectionTypeTXP SectionTypes = [4]byte{'t', 'x', 'p', '1'}
	SectionTypeTXL SectionTypes = [4]byte{'t', 'x', 'l', '1'}
	SectionTypeTXT SectionTypes = [4]byte{'t', 'x', 't', '1'}
	SectionTypePAG SectionTypes = [4]byte{'p', 'a', 'g', '1'}
	SectionTypePAN SectionTypes = [4]byte{'p', 'a', 'n', '1'}
	SectionTypePAS SectionTypes = [4]byte{'p', 'a', 's', '1'}
	SectionTypePAE SectionTypes = [4]byte{'p', 'a', 'e', '1'}
	SectionTypeWND SectionTypes = [4]byte{'w', 'n', 'd', '1'}
	SectionTypePIC SectionTypes = [4]byte{'p', 'i', 'c', '1'}
)

type Writer struct {
	root Root
	*bytes.Buffer
}

func main() {
	if len(os.Args) != 4 {
		log.Println("Usage: NTPGParser [toXML|toNTPG] <input> <output>")
		os.Exit(1)
	}

	action := os.Args[1]
	input := os.Args[2]
	output := os.Args[3]

	switch action {
	case "toXML":
		contents, err := os.ReadFile(input)
		if err != nil {
			panic(err)
		}

		data := ParseNTPG(contents)
		err = os.WriteFile(output, data, 0666)
		if err != nil {
			panic(err)
		}
	case "toNTPG":
		file, err := os.ReadFile(input)
		if err != nil {
			return
		}

		data := WriteNTPG(file)
		err = os.WriteFile(output, data, 0666)
		if err != nil {
			panic(err)
		}
	}

}

func ParseNTPG(input []byte) []byte {
	// Create a new reader
	readable := bytes.NewReader(input)

	var header Header
	err := binary.Read(readable, binary.LittleEndian, &header)
	if err != nil {
		panic(err)
	}

	if !bytes.Equal(headerMagic[:], header.Magic[:]) {
		panic(ErrInvalidFileMagic)
	}

	if readable.Size() != int64(header.FileSize) {
		panic(ErrFileSizeMismatch)
	}

	root := Root{
		XMLName: xml.Name{},
		NAP:     NAPNode{},
	}

	for count := header.SectionCount; count != 0; count-- {
		var sectionHeader SectionHeader
		err = binary.Read(readable, binary.LittleEndian, &sectionHeader)
		if err != nil {
			panic(err)
		}

		// Subtract the header size
		sectionSize := int(sectionHeader.Size) - 8

		if readable.Len() == 0 {
			// If our type is one of the section ending types, we can write then finish.
			switch sectionHeader.Type {
			case SectionTypePAE:
				root.ParsePAE()
			}
			continue
		}

		temp := make([]byte, sectionSize)
		_, err = readable.Read(temp)
		if err != nil {
			panic(err)
		}

		switch sectionHeader.Type {
		case SectionTypeNAP:
			root.ParseNAP(temp)
		case SectionTypeTXP:
			root.ParseTXP(temp)
		case SectionTypeTXL:
			root.ParseTXL(temp)
		case SectionTypeTXT:
			root.ParseTXT(temp)
		case SectionTypePAG:
			root.ParsePage(temp)
		case SectionTypePAN:
			root.ParsePane(temp)
		case SectionTypePAS:
			root.ParsePAS()
		case SectionTypePAE:
			root.ParsePAE()
		case SectionTypeWND:
			root.ParseWND(temp)
		case SectionTypePIC:
			root.ParsePic(temp)
		default:
			break
		}
	}

	data, err := xml.MarshalIndent(root, "", "\t")
	if err != nil {
		panic(err)
	}

	return data
}

func WriteNTPG(data []byte) []byte {
	var root Root
	err := xml.Unmarshal(data, &root)
	if err != nil {
		panic(err)
	}

	header := Header{
		Magic:        headerMagic,
		BOM:          binary.LittleEndian.Uint32([]byte{0xFF, 0xFE, 0x00, 0x02}),
		FileSize:     0,
		HeaderLen:    16,
		SectionCount: 0,
	}

	writer := Writer{
		root:   root,
		Buffer: new(bytes.Buffer),
	}

	writer.write(header)
	writer.WriteNAP()
	writer.WriteTXP()
	writer.WriteTXL()
	writer.WritePage()
	writer.WritePane()

	sectionCount := 5

	for _, pane := range root.Panes {
		if pane.TXT != nil {
			writer.WriteTXT(pane.TXT)
		}

		if pane.PAS != nil {
			writer.WritePAS()
		}

		if pane.PAE != nil {
			writer.WritePAE()
		}

		if pane.WND != nil {
			writer.WriteWND(pane.WND)
		}

		if pane.PIC != nil {
			writer.WritePIC(pane.PIC)
		}

		sectionCount++
	}

	binary.LittleEndian.PutUint32(writer.Bytes()[8:12], uint32(writer.Len()))
	binary.LittleEndian.PutUint16(writer.Bytes()[14:16], uint16(sectionCount))

	return writer.Bytes()
}

func (w *Writer) write(data interface{}) {
	err := binary.Write(w, binary.LittleEndian, data)
	if err != nil {
		panic(err)
	}
}
