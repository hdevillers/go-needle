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

type Search struct {
	RefSeq   []seq.Seq
	QuerySeq []seq.Seq
	NRef     int
	NQuery   int
	Param    needle.Param
	Kept     []bool
	MinSim   float64
}

func newSearch(r, q, f string) *Search {
	var s Search
	s.NRef = loadSeq(r, f, &s.RefSeq)
	s.NQuery = loadSeq(q, f, &s.QuerySeq)
	s.Param = *needle.NewParam()
	fmt.Printf("Read %d sequences in reference sequence file (%s).\n", s.NRef, r)
	fmt.Printf("Read %d sequences in reference sequence file (%s).\n", s.NQuery, q)
	s.Kept = make([]bool, s.NQuery)

	return (&s)
}

func (s *Search) runSearchThread(queryChan chan int, threadChan chan int) {
	for iquery := range queryChan {
		query := s.QuerySeq[iquery]
	rtest:
		for _, seq := range s.RefSeq {
			nw := needle.NewNeedle(query, seq)
			nw.SetParam(&s.Param)

			err := nw.Align()
			if err != nil {
				panic(err)
			}

			if nw.Rst.GetSimilarityPct() >= s.MinSim {
				s.Kept[iquery] = true
				break rtest
			}

		}
	}
	// End of thread process
	threadChan <- 1
}

func (s *Search) save(o string) {
	wkept := seqio.NewWriter(o+"_kept.fasta", "fasta")
	wkept.CheckPanic()
	defer wkept.Close()
	wdisc := seqio.NewWriter(o+"_discarded.fasta", "fasta")
	wdisc.CheckPanic()
	defer wdisc.Close()

	nkept := 0
	ndisc := 0
	for i := range s.Kept {
		if (*&s.Kept)[i] {
			nkept++
			wkept.Write(s.QuerySeq[i])
		} else {
			ndisc++
			wdisc.Write(s.QuerySeq[i])
		}
	}

	fmt.Printf("Kept %d sequences from the query sequences (%d discarded).\n", nkept, ndisc)
}

func main() {
	refInput := flag.String("ref", "", "Input reference sequence file.")
	queryInput := flag.String("query", "", "Input query sequence file.")
	output := flag.String("output", "SeqSearch", "Output file name base.")
	minSim := flag.Float64("min-sim", 80.0, "Minimal similarity threshold.")
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

	// Init. a new search object
	src := newSearch(*refInput, *queryInput, *format)
	src.MinSim = *minSim

	// Reset parameter if needed
	if *gapopen != -1.0 {
		src.Param.SetGapOpen(*gapopen)
	}
	if *gapextend != -1.0 {
		src.Param.SetGapExtend(*gapextend)
	}
	if *endopen != -1.0 {
		src.Param.SetEndOpen(*endopen)
	}
	if *endextend != -1.0 {
		src.Param.SetGapExtend(*endextend)
	}
	if *endweight {
		src.Param.SetEndWeight(true)
	}

	// Prepare chanels
	queryChan := make(chan int)
	threadChan := make(chan int)

	// Launch parallel threaded routines
	for i := 0; i < *threads; i++ {
		go src.runSearchThread(queryChan, threadChan)
	}

	// Feed the query channel
	for i := 0; i < src.NQuery; i++ {
		queryChan <- i
	}
	close(queryChan)

	// Wait threads
	for i := 0; i < *threads; i++ {
		<-threadChan
	}

	// Write out selected sequences
	src.save(*output)
}
