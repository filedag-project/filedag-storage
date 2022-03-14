package diskv

import "testing"

type tpair struct {
	In  string
	Out [2]string
}

var testData []tpair = []tpair{
	{In: "QmbuSHXN9RANoN46sGYdxHyG6SEEcvdEJfXuXfD6EtfTUw", Out: [2]string{"Uw/fT/Et", "Uw/fT/Et/QmbuSHXN9RANoN46sGYdxHyG6SEEcvdEJfXuXfD6EtfTUw"}},
	{In: "QmWRzjAJjRHLsMaWZcpVysHqs22P5eBDpj2G4rYWFDrEKh", Out: [2]string{"Kh/rE/FD", "Kh/rE/FD/QmWRzjAJjRHLsMaWZcpVysHqs22P5eBDpj2G4rYWFDrEKh"}},
	{In: "Qmepbk8EMnA7ssi1vd7A9qUDknogAVWLd8Kk3XqcApk5G5", Out: [2]string{"G5/k5/Ap", "G5/k5/Ap/Qmepbk8EMnA7ssi1vd7A9qUDknogAVWLd8Kk3XqcApk5G5"}},
	{In: "QmT8iNCN13gs2x2563pvv5mexSE5FUVR5gcApfFTvyDMUJ", Out: [2]string{"UJ/DM/vy", "UJ/DM/vy/QmT8iNCN13gs2x2563pvv5mexSE5FUVR5gcApfFTvyDMUJ"}},
}

func TestDefaultShardFun(t *testing.T) {
	for _, td := range testData {
		pp, p, err := DefaultShardFun(td.In)
		if err != nil {
			t.Fatal(err)
		}
		if pp != td.Out[0] || p != td.Out[1] {
			t.Fatal("unmatched output")
		}
	}
}
