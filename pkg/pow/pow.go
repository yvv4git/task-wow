package pow

type Proof struct {
	Hash  []byte `json:"hash"`
	Nonce int    `json:"nonce"`
}

type POW interface {
	SolveChallenge(data []byte) Proof
	Check(data []byte, proof Proof) bool
}
