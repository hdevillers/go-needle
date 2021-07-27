package needle

import (
	"errors"
	"regexp"
	"strconv"
)

type Result struct {
	length   int
	nident   int
	nsimil   int
	ngap     int
	pctident float64
	pctsimil float64
	pctgap   float64
	score    float64
}

type Reader interface {
	ReadString(byte) (string, error)
}

func NewResult() *Result {
	return &Result{}
}

func (r *Result) Parse(out Reader) error {
	var line string
	var err error

	// Init. regex
	reLength := regexp.MustCompile(`^\# Length\: (\d+)`)
	reIdent := regexp.MustCompile(`Identity\:\s+(\d+)`)
	reSimil := regexp.MustCompile(`Similarity\:\s+(\d+)`)
	reGap := regexp.MustCompile(`Gaps\:\s+(\d+)`)
	reScore := regexp.MustCompile(`Score\:\s+([\d+\.]+)`)

	// Skip first lines and read the length
	for err == nil {
		line, err = out.ReadString('\n')
		if reLength.MatchString(line) {
			r.length, err = strconv.Atoi(reLength.FindStringSubmatch(line)[1])
			break
		}
	}
	if err != nil {
		// Even if err is io.EOF
		return errors.New("[Needle parser]: Failed to find alignment length. Please check output format.")
	}

	// Continue and read the identity
	line, err = out.ReadString('\n')
	if err == nil {
		r.nident, err = strconv.Atoi(reIdent.FindStringSubmatch(line)[1])
	}
	if err != nil {
		return errors.New("[Needle parser]: Failed to find identity score. Please check output format.")
	}

	// Continue and read the similarity
	line, err = out.ReadString('\n')
	if err == nil {
		r.nsimil, err = strconv.Atoi(reSimil.FindStringSubmatch(line)[1])
	}
	if err != nil {
		return errors.New("[Needle parser]: Failed to find identity score. Please check output format.")
	}

	// Continue and read the number of gaps
	line, err = out.ReadString('\n')
	if err == nil {
		r.ngap, err = strconv.Atoi(reGap.FindStringSubmatch(line)[1])
	}
	if err != nil {
		return errors.New("[Needle parser]: Failed to find the number of gaps. Please check output format.")
	}

	// Continue and read the score
	line, err = out.ReadString('\n')
	if err == nil {
		r.score, err = strconv.ParseFloat(reScore.FindStringSubmatch(line)[1], 64)
	}
	if err != nil {
		return errors.New("[Needle parser]: Failed to find the alignment score. Please check output format.")
	}

	// Compute ratio
	if r.length > 0 {
		r.pctident = float64(r.nident) / float64(r.length)
		r.pctsimil = float64(r.nsimil) / float64(r.length)
		r.pctgap = float64(r.ngap) / float64(r.length)
	} else {
		return errors.New("[Needle parser]: Alignment length is null.")
	}

	// Catch the aligned sequences (skip line while # is the first character)
	for line[0] == '#' {
		line, err = out.ReadString('\n')
		if err != nil {
			return errors.New("[Needle parser]: ...")
		}
	}

	return nil
}

// GETTERS
func (r *Result) GetSimilarityPct() float64 {
	return r.pctsimil
}
func (r *Result) GetSimilarityCount() int {
	return r.nsimil
}
func (r *Result) GetIdentityPct() float64 {
	return r.pctident
}
func (r *Result) GetIdentityCount() int {
	return r.nident
}
func (r *Result) GetGapPct() float64 {
	return r.pctgap
}
func (r *Result) GetGapCount() int {
	return r.ngap
}
func (r *Result) GetScore() float64 {
	return r.score
}
func (r *Result) GetLength() int {
	return r.length
}
