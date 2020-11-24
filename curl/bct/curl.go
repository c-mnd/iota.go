// Package bct implements the BCT Curl hashing function computing multiple Curl hashes in parallel.
package bct

import (
	"math/bits"

	"github.com/iotaledger/iota.go/consts"
	"github.com/iotaledger/iota.go/curl"
	"github.com/iotaledger/iota.go/trinary"
)

// MaxBatchSize is the maximum number of Curl hashes that can be computed in one batch.
const MaxBatchSize = bits.UintSize

// Curl is the BCT version of the Curl hashing function.
type Curl struct {
	p, n      [curl.StateSize]uint // main batched state of the hash
	direction curl.SpongeDirection // whether the sponge is absorbing or squeezing
}

// NewCurlP81 returns a new BCT Curl-P-81.
func NewCurlP81() *Curl {
	c := &Curl{
		direction: curl.SpongeAbsorbing,
	}
	return c
}

// Reset the internal state of the BCT Curl instance.
func (c *Curl) Reset() {
	c.p = [curl.StateSize]uint{}
	c.n = [curl.StateSize]uint{}
	c.direction = curl.SpongeAbsorbing
}

// Clone returns a deep copy of the current BCT Curl instance.
func (c *Curl) Clone() *Curl {
	return &Curl{
		p:         c.p,
		n:         c.n,
		direction: c.direction,
	}
}

// Absorb fills the states of the sponge with src; each element of src must have the length tritsCount.
// The value tritsCount has to be a multiple of HashTrinarySize.
func (c *Curl) Absorb(src []trinary.Trits, tritsCount int) error {
	if len(src) < 1 || len(src) > MaxBatchSize {
		return consts.ErrInvalidBatchSize
	}
	if tritsCount%consts.HashTrinarySize != 0 {
		return consts.ErrInvalidTritsLength
	}

	if c.direction != curl.SpongeAbsorbing {
		panic("absorb after squeeze")
	}
	for i := 0; i < tritsCount/consts.HashTrinarySize; i++ {
		for j := range src {
			c.in(src[j][i*consts.HashTrinarySize:], uint(j))
		}
		c.transform()
	}
	return nil
}

// Squeeze squeezes out trits of the given length.
// The value tritsCount has to be a multiple of HashTrinarySize.
func (c *Curl) Squeeze(dst []trinary.Trits, tritsCount int) error {
	if len(dst) < 1 || len(dst) > MaxBatchSize {
		return consts.ErrInvalidBatchSize
	}
	if tritsCount%consts.HashTrinarySize != 0 {
		return consts.ErrInvalidSqueezeLength
	}

	for j := range dst {
		dst[j] = make(trinary.Trits, tritsCount)
	}
	for i := 0; i < tritsCount/consts.HashTrinarySize; i++ {
		// during squeezing, we only transform before each squeeze to avoid unnecessary transforms
		if c.direction == curl.SpongeSqueezing {
			c.transform()
		}
		c.direction = curl.SpongeSqueezing
		for j := range dst {
			c.out(dst[j][i*consts.HashTrinarySize:], uint(j))
		}
	}
	return nil
}

// in sets the idx-th entry of the internal state to src.
func (c *Curl) in(src trinary.Trits, idx uint) {
	// bounds check hint to compiler
	if len(src) < consts.HashTrinarySize {
		panic(consts.ErrInvalidTritsLength)
	}

	s := uint(1) << idx
	u := ^s
	for i := 0; i < consts.HashTrinarySize; i++ {
		switch src[i] {
		case 1:
			c.p[i] |= s
			c.n[i] &= u
		case -1:
			c.p[i] &= u
			c.n[i] |= s
		default:
			c.p[i] &= u
			c.n[i] &= u
		}
	}
}

// out extracts the idx-th entry of the internal state to dst.
func (c *Curl) out(dst trinary.Trits, idx uint) {
	// bounds check hint to compiler
	if len(dst) < consts.HashTrinarySize {
		panic(consts.ErrInvalidTritsLength)
	}

	for i := 0; i < consts.HashTrinarySize; i++ {
		p := (c.p[i] >> idx) & 1
		n := (c.n[i] >> idx) & 1

		switch {
		case p != 0:
			dst[i] = 1
		case n != 0:
			dst[i] = -1
		default:
			dst[i] = 0
		}
	}
}

// transform transforms the sponge.
func (c *Curl) transform() {
	var p2, n2 [curl.StateSize]uint
	transformGeneric(&p2, &n2, &c.p, &c.n, curl.NumRounds)
	copy(c.p[:], p2[:])
	copy(c.n[:], n2[:])
}
