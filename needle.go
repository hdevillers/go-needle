package needle

import (
	"github.com/hdevillers/go-seq/seq"
)

const (
	D_GAPOPEN   float64 = 10.0
	D_GAPEXTEND float64 = 0.5
	D_ENDOPEN   float64 = 10.0
	D_ENDEXTEND float64 = 0.5
	D_ENDWEIGHT bool    = false
)

type Needle struct {
	seqa seq.Seq
	seqb seq.Seq
}

func NewNeedle(sa seq.Seq, sb seq.Seq, gapo float64, gape float64, endo float64, ende float64, endw bool) &Needle {
	
}