package main

func (r *Root) ParseTXL(data []byte) {
	// There isn't really any useful data here, it is literally just the nap table but without the strings.
	// First offset is the texture index in u32.
}

func (w *Writer) WriteTXL() {
	header := SectionHeader{
		Type: SectionTypeTXL,
		Size: 12 + uint32(len(w.root.NAP.Contents)*4),
	}

	w.write(header)
	w.write(uint32(len(w.root.NAP.Contents)))
	for i := 0; i < len(w.root.NAP.Contents); i++ {
		w.write(uint32(i))
	}
}
