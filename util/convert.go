package util

func IntToByteArray(n int) [12]byte {
	var byteArray [12]byte

	// Use bitwise AND operations to extract individual bytes
	byteArray[0] = byte(n >> 88)
	byteArray[1] = byte(n >> 80)
	byteArray[2] = byte(n >> 72)
	byteArray[3] = byte(n >> 64)
	byteArray[4] = byte(n >> 56)
	byteArray[5] = byte(n >> 48)
	byteArray[6] = byte(n >> 40)
	byteArray[7] = byte(n >> 32)
	byteArray[8] = byte(n >> 24)
	byteArray[9] = byte(n >> 16)
	byteArray[10] = byte(n >> 8)
	byteArray[11] = byte(n)

	return byteArray
}

func ByteArrayToInt(byteArray [12]byte) int {
	// Convert the [12]byte to a 32-bit integer (int32)
	var intValue int
	for i := 0; i < 4; i++ {
		intValue = (intValue << 8) | int(byteArray[i])
	}
	return intValue
}
