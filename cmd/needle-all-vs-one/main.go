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

type AllOne struct {
	Seqs  []seq.Seq
	Ref   seq.Seq
	Nseqs int
	Param needle.Param
	Simil []float64
	Ident []float64
}

func newAllOne(q, r, f string) *AllOne {
	var ao AllOne

	// Load the queries (All)
	ao.Nseqs = loadSeq(q, f, &ao.Seqs)
	if ao.Nseqs == 0 {
		panic("No sequence found in input query.")
	}

	// Load the reference (One)
	reader := seqio.NewReader(r, f, false) // We suppose the ref not zipped
	reader.CheckPanic()
	defer reader.Close()
	if reader.Next() {
		reader.CheckPanic()
		ao.Ref = reader.Seq()
	} else {
		panic("Failed to read reference sequence.")
	}

	// Init a new parameter setting
	ao.Param = *needle.NewParam()

	ao.Simil = make([]float64, ao.Nseqs)
	ao.Ident = make([]float64, ao.Nseqs)

	return &ao
}

func (ao *AllOne) comparisonThread(seqChan chan int, threadChan chan int) {
	for iSeq := range seqChan {
		seqa := ao.Ref
		seqb := ao.Seqs[iSeq]

		nw := needle.NewNeedle(seqa, seqb)
		nw.SetParam(&ao.Param)

		err := nw.Align()
		if err != nil {
			panic(err)
		}

		ao.Simil[iSeq] = nw.Rst.GetSimilarityPct()
		ao.Ident[iSeq] = nw.Rst.GetIdentityPct()
	}

	// End of thread
	threadChan <- 1
}

func (ao *AllOne) saveResults(o string) {
	// Create output handle and buffer
	f, e := os.Create(o + "_AllvsOne.tsv")
	if e != nil {
		panic(e)
	}
	defer f.Close()
	b := bufio.NewWriter(f)

	// Write out header
	b.WriteString("Query\tRef\tSimilarity\tIdentity\n")

	for i := 0; i < ao.Nseqs; i++ {
		b.WriteString(ao.Seqs[i].Id)
		b.WriteByte('\t')
		b.WriteString(ao.Ref.Id)
		b.WriteByte('\t')
		b.WriteString(fmt.Sprintf("%0.3f\t%0.3f\n", ao.Simil[i], ao.Ident[i]))
	}
	b.Flush()
}

func main() {
	ref := flag.String("ref", "", "Reference input sequence file.")
	query := flag.String("query", "", "Query input sequence file.")
	output := flag.String("output", "AllVsAll", "Output file name base.")
	format := flag.String("format", "fasta", "Sequence file format.")
	gapopen := flag.Float64("gapopen", -1.0, "Gap open penality.")
	gapextend := flag.Float64("gapextend", -1.0, "Gap extend penality.")
	endopen := flag.Float64("endopen", -1.0, "End gap open penality.")
	endextend := flag.Float64("endextend", -1.0, "End gap extend penality.")
	endweight := flag.Bool("endweight", false, "Apply end gap penality.")
	threads := flag.Int("threads", 2, "Number of threads.")
	flag.Parse()

	if *ref == "" {
		panic("You must provide a reference input sequence file.")
	}
	if *query == "" {
		panic("You must provide a query input sequence file.")
	}

	// Init. an AllAll object
	ao := newAllOne(*query, *ref, *format)

	// Reset parameter if needed
	if *gapopen != -1.0 {
		ao.Param.SetGapOpen(*gapopen)
	}
	if *gapextend != -1.0 {
		ao.Param.SetGapExtend(*gapextend)
	}
	if *endopen != -1.0 {
		ao.Param.SetEndOpen(*endopen)
	}
	if *endextend != -1.0 {
		ao.Param.SetGapExtend(*endextend)
	}
	if *endweight {
		ao.Param.SetEndWeight(true)
	}

	// Prepare chanels
	seqChan := make(chan int)
	threadChan := make(chan int)

	// Launch parallel comparison routine
	for i := 0; i < *threads; i++ {
		go ao.comparisonThread(seqChan, threadChan)
	}

	// Send combs to channel
	for i := 0; i < ao.Nseqs; i++ {
		seqChan <- i
	}
	close(seqChan)

	// Wait for all threads
	for i := 0; i < *threads; i++ {
		<-threadChan
	}

	// Write results
	ao.saveResults(*output)
}
