package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type lineType int

const (
	lineTypeUnknown lineType = iota
	lineTypeStartTest
	lineTypeEndTest
	lineTypeComment
)

type parsedLine struct {
	lineType lineType
	line     string
}

var (
	rgxStartTestSingleComment = regexp.MustCompile(`^\s*(--)*\s*START_TEST(.*)`)
	rgxStartTestGoupComment   = regexp.MustCompile(`^\s*\/\*\s*START_TEST(.*)`)
	rgxEndTestSingleComment   = regexp.MustCompile(`^\s*(--)*\s*END_TEST(.*)`)
	rgxEndTestGroupComment    = regexp.MustCompile(`(.*)END_TEST$`)
	rgxEndGroupComment        = regexp.MustCompile(`\s*\*\/`)
	rgxStartSingleComment     = regexp.MustCompile(`^\s*--+(.*)`)
)

func parseLine(line string) parsedLine {
	pl := parsedLine{
		line: line,
	}

	switch {
	case rgxStartTestSingleComment.MatchString(line), rgxStartTestGoupComment.MatchString(line):
		pl.lineType = lineTypeStartTest
	case rgxEndTestSingleComment.MatchString(line), rgxEndTestGroupComment.MatchString(line):
		pl.lineType = lineTypeEndTest
	case rgxEndGroupComment.MatchString(line):
		pl.lineType = lineTypeComment
	case rgxStartSingleComment.MatchString(line):
		pl.lineType = lineTypeComment
	default:
		pl.lineType = lineTypeUnknown
	}

	return pl
}

func parseFile(filename string) ([]parsedLine, error) {
	rdr, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}

	bufioScanner := bufio.NewScanner(rdr)

	res := []parsedLine{}

	for bufioScanner.Scan() {
		line := bufioScanner.Text()
		res = append(res, parseLine(line))
	}

	return res, nil
}
