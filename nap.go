package main

import (
	"bytes"
	"encoding/binary"
)

func (r *Root) ParseNAP(data []byte) {
	reader := bytes.NewReader(data)
	var numberOfNAPs uint32
	err := binary.Read(reader, binary.LittleEndian, &numberOfNAPs)
	if err != nil {
		panic(err)
	}

	offsets := make([]uint32, numberOfNAPs+1)
	for i := uint32(0); i < numberOfNAPs+1; i++ {
		err = binary.Read(reader, binary.LittleEndian, &offsets[i])
		if err != nil {
			panic(err)
		}
	}

	// Read the strings now.
	strings := make([]string, numberOfNAPs)
	for i := uint32(0); i < numberOfNAPs; i++ {
		size := offsets[i+1] - offsets[i]
		str := make([]byte, size)
		reader.Read(str)

		strings[i] = string(str[:size-1])
	}

	r.NAP = NAPNode{
		Contents: strings,
	}
}

func (w *Writer) WriteNAP() {
	header := SectionHeader{
		Type: SectionTypeNAP,
		Size: 0,
	}

	size := uint32(12 + (len(w.root.NAP.Contents)+1)*4)
	for _, content := range w.root.NAP.Contents {
		size += uint32(len(content) + 1)
	}

	header.Size = size

	offsets := make([]uint32, len(w.root.NAP.Contents)+1)

	// First offset is relative to the number of offsets.
	offsets[0] = uint32((len(w.root.NAP.Contents) + 1) * 4)
	for i := 0; i < len(w.root.NAP.Contents); i++ {
		// Add extra byte for null terminator
		offsets[i+1] = offsets[i] + uint32(len(w.root.NAP.Contents[i])) + 1
	}

	w.write(header)
	w.write(uint32(len(w.root.NAP.Contents)))
	w.write(offsets)

	// Write strings
	for _, content := range w.root.NAP.Contents {
		w.WriteString(content)
		w.WriteByte(0)
	}
}
