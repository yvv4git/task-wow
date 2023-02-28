package pow_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yvv4git/task-wow/pkg/pow"
)

func TestSHA256_Generate(t *testing.T) {
	type args struct {
		data       []byte
		difficulty int
	}

	tests := []struct {
		name string
		args args
		want pow.Proof
		desc string
	}{
		{
			name: "CASE-1",
			args: args{
				data:       []byte("admin@gmail.com"),
				difficulty: 1,
			},
			want: pow.Proof{
				Hash: []byte{
					0x0, 0x9, 0xf5, 0x18, 0x3e, 0xc0, 0xf5, 0x1b, 0xc5, 0xa, 0x9, 0x3c, 0x1a, 0x16, 0x48, 0x67, 0x66,
					0xbf, 0xca, 0x74, 0xa6, 0xe5, 0x52, 0xf9, 0xb3, 0x16, 0x22, 0x52, 0xc, 0x36, 0x4a, 0xe5,
				},
				Nonce: 87,
			},
			desc: "Very fast ~ 0.00s",
		},
		{
			name: "CASE-2",
			args: args{
				data:       []byte("admin@gmail.com"),
				difficulty: 2,
			},
			want: pow.Proof{
				Hash: []byte{
					0x0, 0x0, 0x71, 0xff, 0x79, 0xd0, 0xf1, 0xb4, 0x5, 0x87, 0x9b, 0xf7, 0xfe, 0xed, 0xdc, 0x55, 0xac,
					0x8f, 0xb, 0xe1, 0xa6, 0xb3, 0x32, 0xfa, 0xf6, 0x14, 0x19, 0xf7, 0xa, 0xb7, 0x75, 0xe9,
				},
				Nonce: 98045,
			},
			desc: "Also very fast ~ 0.04s",
		},
		{
			name: "CASE-3",
			args: args{
				data:       []byte("admin@gmail.com"),
				difficulty: 3,
			},
			want: pow.Proof{
				Hash: []byte{
					0x0, 0x0, 0x0, 0x79, 0x3a, 0xb2, 0x1f, 0x2, 0x3e, 0x14, 0xd2, 0x2a, 0x45, 0x34, 0xd3, 0x42, 0x92,
					0xd2, 0x3, 0xb7, 0x41, 0xfa, 0x5e, 0xb, 0x23, 0x7b, 0x97, 0x93, 0xf0, 0xd, 0xf9, 0x79,
				},
				Nonce: 5675647,
			},
			desc: "We have to wait >1s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			algorithm := pow.NewSHA256(tt.args.difficulty)
			proof := algorithm.SolveChallenge(tt.args.data)
			assert.Equal(t, tt.want, proof)
			assert.True(t, algorithm.Check(tt.args.data, proof))
		})
	}
}

func BenchmarkSHA256_Generate(b *testing.B) {
	benchmarks := []struct {
		name       string
		data       []byte
		difficulty int
	}{
		{
			name:       "CASE-1",
			data:       []byte("admin@gmail.com"),
			difficulty: 1,
		},
		{
			name:       "CASE-2",
			data:       []byte("admin@gmail.com"),
			difficulty: 2,
		},
		{
			name:       "CASE-3",
			data:       []byte("admin@gmail.com"),
			difficulty: 3,
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				algorithm := pow.NewSHA256(bm.difficulty)
				algorithm.SolveChallenge(bm.data)
			}
		})
	}
}
