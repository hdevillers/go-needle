package main

import (
	"flag"
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
	gapopen := flag.Float64("gapopen", needle.D_GAPOPEN, "Gap open penality.")
	gapextend := flag.Float64("gapextend", needle.D_GAPEXTEND, "Gap extend penality.")
	endopen := flag.Float64("endopen", needle.D_ENDOPEN, "End gap open penality.")
	endextend := flag.Float64("endextend", needle.D_ENDEXTEND, "End gap extend penality.")
	endweight := flag.Bool("endweight", needle.D_ENDWEIGHT, "Apply end gap penality.")

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

}
