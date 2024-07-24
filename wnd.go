package main

import (
	"bytes"
	"encoding/binary"
)

type WND struct {
	XTranslation      uint16 `xml:"x_translation"`
	YTranslation      uint16 `xml:"y_translation"`
	Width             uint16 `xml:"width"`
	Height            uint16 `xml:"height"`
	SomeValue         uint16 `xml:"some_value"`
	NumberOfMaterials uint16 `xml:"number_of_materials"`
}

type Material struct {
	Unk   uint16 `xml:"unk"`
	Value uint16 `xml:"material_index"`
}

func (r *Root) ParseWND(data []byte) {
	read := bytes.NewReader(data)
	var wnd WND
	err := binary.Read(read, binary.LittleEndian, &wnd)
	if err != nil {
		panic(err)
	}

	mats := make([]Material, wnd.NumberOfMaterials)
	for i := uint16(0); i < wnd.NumberOfMaterials; i++ {
		err := binary.Read(read, binary.LittleEndian, &mats[i])
		if err != nil {
			panic(err)
		}
	}

	r.Panes = append(r.Panes, Children{WND: &XMLWND{
		WND:       wnd,
		Materials: mats,
	}})
}

func (w *Writer) WriteWND(wnd *XMLWND) {
	header := SectionHeader{
		Type: SectionTypeWND,
		Size: uint32(20 + (4 * wnd.NumberOfMaterials)),
	}

	w.write(header)
	w.write(wnd.WND)

	for _, material := range wnd.Materials {
		w.write(material)
	}
}
