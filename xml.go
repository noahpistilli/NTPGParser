package main

import "encoding/xml"

type Root struct {
	XMLName xml.Name   `xml:"root"`
	TXP     TXPNode    `xml:"txp1"`
	NAP     NAPNode    `xml:"nap1"`
	PAG     Dimensions `xml:"pag1"`
	PAN     Dimensions `xml:"pan1"`
	Panes   []Children `xml:"children"`
	txp     [][]uint16
}

type NAPNode struct {
	Contents []string `xml:"contents"`
}

type TXPNode struct {
	Content string `xml:"content"`
}

type TXTNode struct {
	XTranslation uint16        `xml:"x_translation"`
	YTranslation uint16        `xml:"y_translation"`
	Width        uint16        `xml:"width"`
	Height       uint16        `xml:"height"`
	TextIndex    uint16        `xml:"text_index"`
	XScale       uint16        `xml:"x_scale"`
	YScale       uint16        `xml:"y_scale"`
	SomeValue    uint32        `xml:"somevalue"`
	Lines        []StringLines `xml:"lines"`
}

type StringLines struct {
	Width  uint16 `xml:"width"`
	Height uint16 `xml:"height"`
	String string `xml:"string"`
}

type Dimensions struct {
	Width  uint16 `xml:"width"`
	Height uint16 `xml:"height"`
}

// Children contains all the possible children a brlyt can contain.
// This is needed for unmarshalling when we put together a new brlyt.
type Children struct {
	PAS *XMLPAS  `xml:"pas1"`
	TXT *TXTNode `xml:"txt1"`
	PAE *XMLPAE  `xml:"pae1"`
	WND *XMLWND  `xml:"wnd1"`
	PIC *PIC     `xml:"pic1"`
}

type XMLPAS struct{}
type XMLPAE struct{}

type XMLWND struct {
	WND
	Materials []Material `xml:"materials"`
}
