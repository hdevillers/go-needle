package main

import (
	"flag"
	"fmt"
	"regexp"

	"github.com/hdevillers/go-needle"
	"github.com/hdevillers/go-seq/seq"
	"github.com/hdevillers/go-seq/seqio"
)

func loadSeq(i, f string, a *[]seq.Seq) int {
	nseq := 0
	isgz := false
	if regexp.MustCompile(`\.gz$`).MatchString(i) {
		isgz = true
	}
	reader := seqio.NewReader(i, f, isgz)
	reader.CheckPanic()
	defer reader.Close()

	for reader.Next() {
		reader.CheckPanic()
		*a = append(*a, reader.Seq())
		nseq++
	}

	return nseq
}

func runSearchThread(queryChan chan int, querySeq *[]seq.Seq, refSeq *[]seq.Seq, p *needle.Param, kept *[]int) {
	for iquery := range queryChan {
		query := (*querySeq)[iquery]
		for _, seq := range *refSeq {
			nw := needle.NewNeedle(query, seq)
			nw.SetParam(p)

			err := nw.Align()
			if err != nil {
				panic(err)
			}

		}
	}
}

func main() {
	refInput := flag.String("ref", "", "Input reference sequence file.")
	queryInput := flag.String("query", "", "Input query sequence file.")
	minSimil := flag.Float64("min-sim", 80.0, "Minimal similarity threshold.")
	format := flag.String("format", "fasta", "Sequence file format.")
	gapopen := flag.Float64("gapopen", -1.0, "Gap open penality.")
	gapextend := flag.Float64("gapextend", -1.0, "Gap extend penality.")
	endopen := flag.Float64("endopen", -1.0, "End gap open penality.")
	endextend := flag.Float64("endextend", -1.0, "End gap extend penality.")
	endweight := flag.Bool("endweight", false, "Apply end gap penality.")
	threads := flag.Int("threads", 4, "Number of threads.")
	flag.Parse()

	// Check arguments
	if *refInput == "" {
		panic("You must provide a reference sequence file.")
	}
	if *queryInput == "" {
		panic("You must provide a query sequence file.")
	}

	// Prepare needle parameter setting
	param := needle.NewParam()
	if *gapopen != -1.0 {
		param.SetGapOpen(*gapopen)
	}
	if *gapextend != -1.0 {
		param.SetGapExtend(*gapextend)
	}
	if *endopen != -1.0 {
		param.SetEndOpen(*endopen)
	}
	if *endextend != -1.0 {
		param.SetGapExtend(*endextend)
	}
	if *endweight {
		param.SetEndWeight(true)
	}

	// Load sequences
	var refSeq []seq.Seq
	nref := loadSeq(*refInput, *format, &refSeq)
	fmt.Printf("Read %d sequences in reference sequence file (%s).", nref, *refInput)
	var querySeq []seq.Seq
	nquery := loadSeq(*queryInput, *format, &querySeq)
	fmt.Printf("Read %d sequences in reference sequence file (%s).", nquery, *queryInput)

	// Prepare chanels
	queryChan := make(chan int)

}
