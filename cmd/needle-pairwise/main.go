package main

import (
	"flag"
	"fmt"
	"regexp"

	"github.com/hdevillers/go-needle"
	"github.com/hdevillers/go-seq/seq"
	"github.com/hdevillers/go-seq/seqio"
)

func loadSeq(i, f string) seq.Seq {
	isgz := false
	if regexp.MustCompile(`\.gz$`).MatchString(i) {
		isgz = true
	}
	reader := seqio.NewReader(i, f, isgz)
	reader.CheckPanic()
	defer reader.Close()
	reader.Next()
	reader.CheckPanic()
	seq := reader.Seq() // Get only the first sequence on the file
	return seq
}

func main() {
	seqa := flag.String("seqa", "", "First sequence file.")
	seqb := flag.String("seqb", "", "Second sequence file.")
	format := flag.String("format", "fasta", "Sequence format.")
	gapopen := flag.Float64("gapopen", -1.0, "Gap open penality.")
	gapextend := flag.Float64("gapextend", -1.0, "Gap extend penality.")
	endopen := flag.Float64("endopen", -1.0, "End gap open penality.")
	endextend := flag.Float64("endextend", -1.0, "End gap extend penality.")
	endweight := flag.Bool("endweight", false, "Apply end gap penality.")
	flag.Parse()

	// Check input arguments
	if *seqa == "" {
		panic("You must provide a file name for the first sequence.")
	}
	if *seqb == "" {
		panic("You must provide a file name for the second sequence.")
	}

	// Load each sequences
	sa := loadSeq(*seqa, *format)
	sb := loadSeq(*seqb, *format)

	// Create the Needleman Wunch object
	nw := needle.NewNeedle(sa, sb)

	// Apply user parameter if required
	if *gapopen != -1.0 {
		nw.Par.SetGapOpen(*gapopen)
	}
	if *gapextend != -1.0 {
		nw.Par.SetGapExtend(*gapextend)
	}
	if *endopen != -1.0 {
		nw.Par.SetEndOpen(*endopen)
	}
	if *endextend != -1.0 {
		nw.Par.SetEndExtend(*endextend)
	}
	if *endweight {
		nw.Par.SetEndWeight(*endweight)
	}

	err := nw.Align()
	if err != nil {
		panic(err)
	}

	fmt.Println("Similarity:", nw.Rst.GetIdentityPct())
}
