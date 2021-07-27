package needle

import (
	"github.com/hdevillers/go-seq/seq"
)

type Needle struct {
	Sqa seq.Seq
	Sqb seq.Seq
	Par *Param
	Run bool
}

func NewNeedle(sa seq.Seq, sb seq.Seq) *Needle {
	var nw Needle
	nw.Sqa = sa
	nw.Sqb = sb
	nw.Par = NewParam()
	nw.Run = false

	return &nw
}
