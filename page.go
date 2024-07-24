package main

import (
	"bytes"
	"encoding/binary"
)

func (r *Root) ParsePage(data []byte) {
	var page Dimensions
	err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &page)
	if err != nil {
		panic(err)
	}

	r.PAG = page
}

func (w *Writer) WritePage() {
	header := SectionHeader{
		Type: SectionTypePAG,
		Size: 16,
	}

	w.write(header)
	w.write(w.root.PAG)
	w.write(binary.LittleEndian.Uint32([]byte{0xFF, 0x7F, 0x00, 0x00}))
}
