package main

import (
	"bytes"
	"encoding/binary"
)

type PIC struct {
	XTranslation  uint16 `xml:"x_translation"`
	YTranslation  uint16 `xml:"y_translation"`
	Width         uint16 `xml:"width"`
	Height        uint16 `xml:"height"`
	MaterialIndex uint32 `xml:"material_index"`
}

func (r *Root) ParsePic(data []byte) {
	var pic PIC
	err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &pic)
	if err != nil {
		panic(err)
	}

	r.Panes = append(r.Panes, Children{PIC: &pic})
}

func (w *Writer) WritePIC(pic *PIC) {
	header := SectionHeader{
		Type: SectionTypePIC,
		Size: uint32(20),
	}

	w.write(header)
	w.write(pic)
}
