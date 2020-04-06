package shared

func Encrypt(data, prefix, xor []byte) []byte {
	for i := range data {
		var b byte
		if i < 32 {
			b = prefix[i]
		} else {
			b = data[i-32]
		}
		b ^= xor[i%32]
		data[i] ^= b
	}
	return data
}

func Decrypt(data, prefix, xor []byte) []byte {
	for i := len(data) - 1; i > -1; i-- {
		var b byte
		if i < 32 {
			b = prefix[i]
		} else {
			b = data[i-32]
		}
		b ^= xor[i%32]
		data[i] ^= b
	}
	return data
}
