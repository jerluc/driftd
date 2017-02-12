package main

// Turns an array of byte arrays into one
// big byte array
func packBytes(bas... []byte) []byte {
	var packed []byte
	for _, ba := range bas {
		packed = append(packed, ba...)
	}
	return packed
}

// Trims off extra null bytes from the
// end of a byte array. Note that this
// returns only a slice into the original
// array
func trimmed(bs []byte) []byte {
	for i, b := range bs {
		if b == 0x00 {
			return bs[:i]
		}
	}
	return bs
}
