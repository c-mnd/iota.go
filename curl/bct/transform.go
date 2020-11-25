package bct

import (
	"github.com/iotaledger/iota.go/curl"
)

func transformGeneric(pto, nto, pfrom, nfrom *[curl.StateSize]uint, rounds uint) {
	for r := rounds; r > 0; r-- {
		for i := 0; i < curl.StateSize-2; i += 3 {
			t0 := curl.Indices[i+0]
			t1 := curl.Indices[i+1]
			t2 := curl.Indices[i+2]
			t3 := curl.Indices[i+3]

			p0, n0 := pfrom[t0], nfrom[t0]
			p1, n1 := pfrom[t1], nfrom[t1]
			p2, n2 := pfrom[t2], nfrom[t2]
			p3, n3 := pfrom[t3], nfrom[t3]

			pto[i+0], nto[i+0] = sBox(p0, n0, p1, n1)
			pto[i+1], nto[i+1] = sBox(p1, n1, p2, n2)
			pto[i+2], nto[i+2] = sBox(p2, n2, p3, n3)
		}
		// swap buffers
		pfrom, pto = pto, pfrom
		nfrom, nto = nto, nfrom
	}
}

func sBox(xP, xN, yP, yN uint) (uint, uint) {
	tmp := xN ^ yP
	return tmp &^ xP, ^tmp &^ (xP ^ yN)
}
