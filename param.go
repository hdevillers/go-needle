package needle

import (
	"fmt"
	"os/exec"
)

const (
	D_GAPOPEN   float64 = 10.0
	D_GAPEXTEND float64 = 0.5
	D_ENDOPEN   float64 = 10.0
	D_ENDEXTEND float64 = 0.5
	D_ENDWEIGHT bool    = false
)

type Param struct {
	gapopen   float64
	gapextend float64
	endweight bool
	endopen   float64
	endextend float64
}

func NewParam() *Param {
	var p Param
	p.gapopen = D_GAPOPEN
	p.gapextend = D_GAPEXTEND
	p.endweight = D_ENDWEIGHT
	p.endopen = D_ENDOPEN
	p.endextend = D_ENDEXTEND

	return &p
}

// SETTERS
func (p *Param) SetGapOpen(v float64) {
	p.gapopen = v
}
func (p *Param) SetGapExtend(v float64) {
	p.gapextend = v
}
func (p *Param) SetEndOpen(v float64) {
	p.endopen = v
}
func (p *Param) SetEndExtend(v float64) {
	p.endextend = v
}
func (p *Param) SetEndWeight(v bool) {
	p.endweight = v
}

// GETTERS
func (p *Param) GetGapOpen() float64 {
	return p.gapopen
}
func (p *Param) GetGapExtend() float64 {
	return p.gapextend
}
func (p *Param) GetEndOpen() float64 {
	return p.endopen
}
func (p *Param) GetEndExtend() float64 {
	return p.endextend
}
func (p *Param) GetEndWeight() bool {
	return p.endweight
}

// GENERATE THE COMMAND LINE
func (p *Param) GetCmd(sa, sb string) *exec.Cmd {
	if p.endweight {
		return exec.Command(
			"needle",
			"-asequence", fmt.Sprintf("asis::%s", sa),
			"-bsequence", fmt.Sprintf("asis::%s", sb),
			"-gapopen", fmt.Sprintf("%f", p.gapopen),
			"-gapextend", fmt.Sprintf("%f", p.gapextend),
			"-outfile", "stdout",
			"-endweight", "Y",
			"-endopen", fmt.Sprintf("%f", p.endopen),
			"-endextend", fmt.Sprintf("%f", p.endextend),
		)
	} else {
		return exec.Command(
			"needle",
			"-asequence", fmt.Sprintf("asis::%s", sa),
			"-bsequence", fmt.Sprintf("asis::%s", sb),
			"-gapopen", fmt.Sprintf("%f", p.gapopen),
			"-gapextend", fmt.Sprintf("%f", p.gapextend),
			"-outfile", "stdout",
		)
	}
}
