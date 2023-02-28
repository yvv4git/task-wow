package pow

import (
	"bytes"
	"crypto/sha256"
	"strconv"
)

type SHA256 struct {
	maxNonce int
	prefix   []byte
}

func NewSHA256(difficulty int) *SHA256 {
	const defaultMaxNonce = 100_000_000_000

	prefix := make([]byte, difficulty)
	bufferIdx := 0
	for n := 0; n < difficulty; n++ {
		bufferIdx += copy(prefix[bufferIdx:], []byte{0x00})
	}

	return &SHA256{
		maxNonce: defaultMaxNonce,
		prefix:   prefix,
	}
}

func (s *SHA256) SolveChallenge(data []byte) Proof {
	for nonce := 0; nonce < s.maxNonce; nonce++ {
		hash := generateSHA256(data, nonce)

		if bytes.HasPrefix(hash[:], s.prefix) {
			return Proof{
				Hash:  hash[:],
				Nonce: nonce,
			}
		}
	}

	return Proof{}
}

func (s *SHA256) Check(data []byte, proof Proof) bool {
	hash := generateSHA256(data, proof.Nonce)

	return bytes.HasPrefix(hash[:], s.prefix)
}

func generateSHA256(data []byte, nonce int) [32]byte {
	var buffer bytes.Buffer

	buffer.WriteString(strconv.Itoa(nonce))
	buffer.WriteRune(':')
	buffer.Write(data)

	return sha256.Sum256(buffer.Bytes())
}
