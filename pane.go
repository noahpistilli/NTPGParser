package main

import (
	"bytes"
	"encoding/binary"
)

func (r *Root) ParsePane(data []byte) {
	// TODO: is there only the root pane in this file format?

	var pane Dimensions
	err := binary.Read(bytes.NewReader(data[4:]), binary.LittleEndian, &pane)
	if err != nil {
		panic(err)
	}

	r.PAN = pane
}

func (r *Root) ParsePAS() {
	r.Panes = append(r.Panes, Children{
		PAS: &XMLPAS{},
	})
}

func (r *Root) ParsePAE() {
	r.Panes = append(r.Panes, Children{
		PAE: &XMLPAE{},
	})
}

func (w *Writer) WritePane() {
	header := SectionHeader{
		Type: SectionTypePAN,
		Size: 16,
	}

	w.write(header)
	w.write(uint32(0))
	w.write(w.root.PAN)
}

func (w *Writer) WritePAS() {
	header := SectionHeader{
		Type: SectionTypePAS,
		Size: 12,
	}

	w.write(header)
	w.write(uint32(len(w.root.Panes) - 2))
}

func (w *Writer) WritePAE() {
	header := SectionHeader{
		Type: SectionTypePAE,
		Size: 8,
	}

	w.write(header)
}
