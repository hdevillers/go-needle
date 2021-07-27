package needle

import (
	"bytes"

	"github.com/hdevillers/go-seq/seq"
)

type Needle struct {
	Sqa seq.Seq
	Sqb seq.Seq
	Par *Param
	Rst *Result
}

func NewNeedle(sa seq.Seq, sb seq.Seq) *Needle {
	var nw Needle
	nw.Sqa = sa
	nw.Sqb = sb
	nw.Par = NewParam()
	nw.Rst = NewResult()
	return &nw
}

func (nw *Needle) Align() error {
	// Get the command line to execute
	cmd := nw.Par.GetCmd(
		string(nw.Sqa.Sequence),
		string(nw.Sqb.Sequence),
	)

	// Run the command and catch the stdout
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	// Create the result object from the out report
	rout := bytes.NewBuffer(out)
	err = nw.Rst.Parse(rout)

	return err
}
