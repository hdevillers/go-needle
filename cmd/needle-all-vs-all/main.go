package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/hdevillers/go-needle"
	"github.com/hdevillers/go-seq/seq"
	"github.com/hdevillers/go-seq/seqio"
)

const (
	RES_COUNT = 2
	RES_SIMIL = 0
	RES_IDENT = 1
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

type AllAll struct {
	Seqs    []seq.Seq
	Nseqs   int
	Param   needle.Param
	Ncombs  int
	Combs   [][]int
	Restuls [][]float64
}

func newAllAll(i, f string) *AllAll {
	var aa AllAll
	aa.Nseqs = loadSeq(i, f, &aa.Seqs)

	// If the is less than 3 sequences then quit
	if aa.Nseqs <= 2 {
		panic("Input sequence file must contain more than 2 sequences.")
	}

	// Compute the number of combinations
	aa.Ncombs = (aa.Nseqs * (aa.Nseqs - 1)) / 2

	// Init. a new parameter setting
	aa.Param = *needle.NewParam()

	// Init. Combs and Results arrays
	aa.Combs = make([][]int, aa.Ncombs)
	aa.Restuls = make([][]float64, aa.Ncombs)
	iComb := 0
	for i := 0; i < (aa.Nseqs - 1); i++ {
		for j := (i + 1); j < aa.Nseqs; j++ {
			aa.Combs[iComb] = make([]int, 2)
			aa.Combs[iComb][0] = i
			aa.Combs[iComb][1] = j
			aa.Restuls[iComb] = make([]float64, RES_COUNT)
			iComb++
		}
	}

	return (&aa)
}

func (aa *AllAll) comparisonThread(combChan chan int, threadChan chan int) {
	for iComb := range combChan {
		seqa := aa.Seqs[aa.Combs[iComb][0]]
		seqb := aa.Seqs[aa.Combs[iComb][1]]

		nw := needle.NewNeedle(seqa, seqb)
		nw.SetParam(&aa.Param)

		err := nw.Align()
		if err != nil {
			panic(err)
		}

		// Store results
		aa.Restuls[iComb][RES_SIMIL] = nw.Rst.GetSimilarityPct()
		aa.Restuls[iComb][RES_IDENT] = nw.Rst.GetIdentityPct()
	}

	// End of the thread
	threadChan <- 1
}

func (aa *AllAll) savePairs(o string) {
	// Create output handle and buffer
	f, e := os.Create(o + "_pairs.tsv")
	if e != nil {
		panic(e)
	}
	defer f.Close()
	b := bufio.NewWriter(f)

	// Write out header
	b.WriteString("Seqa\tSeqb\tSimilarity\tIdentity\n")

	// Write out each results
	for i := 0; i < aa.Ncombs; i++ {
		ida := aa.Seqs[aa.Combs[i][0]].Id
		idb := aa.Seqs[aa.Combs[i][1]].Id

		b.WriteString(ida)
		b.WriteByte('\t')
		b.WriteString(idb)
		b.WriteByte('\t')
		b.WriteString(fmt.Sprintf("%0.3f\t%0.3f\n", aa.Restuls[i][RES_SIMIL], aa.Restuls[i][RES_IDENT]))
	}
	b.Flush()
}

func main() {
	input := flag.String("input", "", "Input sequence file.")
	output := flag.String("output", "AllVsAll", "Output file name base.")
	format := flag.String("format", "fasta", "Sequence file format.")
	gapopen := flag.Float64("gapopen", -1.0, "Gap open penality.")
	gapextend := flag.Float64("gapextend", -1.0, "Gap extend penality.")
	endopen := flag.Float64("endopen", -1.0, "End gap open penality.")
	endextend := flag.Float64("endextend", -1.0, "End gap extend penality.")
	endweight := flag.Bool("endweight", false, "Apply end gap penality.")
	threads := flag.Int("threads", 2, "Number of threads.")
	flag.Parse()

	if *input == "" {
		panic("You must provide an input sequence file.")
	}

	// Init. an AllAll object
	aa := newAllAll(*input, *format)

	// Reset parameter if needed
	if *gapopen != -1.0 {
		aa.Param.SetGapOpen(*gapopen)
	}
	if *gapextend != -1.0 {
		aa.Param.SetGapExtend(*gapextend)
	}
	if *endopen != -1.0 {
		aa.Param.SetEndOpen(*endopen)
	}
	if *endextend != -1.0 {
		aa.Param.SetGapExtend(*endextend)
	}
	if *endweight {
		aa.Param.SetEndWeight(true)
	}

	// Prepare chanels
	combChan := make(chan int)
	threadChan := make(chan int)

	// Launch parallel comparison routine
	for i := 0; i < *threads; i++ {
		go aa.comparisonThread(combChan, threadChan)
	}

	// Send combs to channel
	for i := 0; i < aa.Ncombs; i++ {
		combChan <- i
	}
	close(combChan)

	// Wait for all threads
	for i := 0; i < *threads; i++ {
		<-threadChan
	}

	// Write results
	aa.savePairs(*output)
}
