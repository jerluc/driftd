package main

func packBytes(bas... []byte) []byte {
	var packed []byte
	for _, ba := range bas {
		packed = append(packed, ba...)
	}
	return packed
}

func trimmed(packet []byte) []byte {
	for i, b := range packet {
		if b == 0x00 {
			return packet[:i]
		}
	}
	return packet
}
