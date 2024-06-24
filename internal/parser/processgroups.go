package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type instructionPrefix int

const (
	instructionPrefixUnknown instructionPrefix = iota
	instructionPrefixStartTest
	instructionPrefixEndTest
	instructionPrefixCount
	instructionPrefixFile
	instructionPrefixRow
)

func (ip instructionPrefix) String() string {
	switch ip {
	case instructionPrefixStartTest:
		return "START_TEST"
	case instructionPrefixEndTest:
		return "END_TEST"
	case instructionPrefixCount:
		return "COUNT"
	case instructionPrefixFile:
		return "FILE"
	case instructionPrefixRow:
		return "ROW"
	default:
		return "UNKNOWN"
	}
}

func buildMapPrefix() map[string]instructionPrefix {
	return map[string]instructionPrefix{
		"START_TEST": instructionPrefixStartTest,
		"END_TEST":   instructionPrefixEndTest,
		"COUNT":      instructionPrefixCount,
		"FILE":       instructionPrefixFile,
		"ROW":        instructionPrefixRow,
	}
}

type prefixAllowance int

const (
	prefixAllowanceUnknown prefixAllowance = iota
	prefixAllowanceSingle
	prefixAllowanceMultiple
)

func (ip instructionPrefix) Allowance() prefixAllowance {
	switch ip {
	case instructionPrefixStartTest:
		return prefixAllowanceSingle
	case instructionPrefixEndTest:
		return prefixAllowanceSingle
	case instructionPrefixCount:
		return prefixAllowanceSingle
	case instructionPrefixFile:
		return prefixAllowanceSingle
	case instructionPrefixRow:
		return prefixAllowanceMultiple
	default:
		return prefixAllowanceUnknown
	}
}

func checkCombinedInstructions(instrs []*outputInstruction) error {
	uniquePrefixes := make(map[instructionPrefix]struct{})

	for _, instr := range instrs {
		uniquePrefixes[instr._type] = struct{}{}
	}

	return checkValidUniquePrefixes(uniquePrefixes)
}

var rgxInstructionPrefix = regexp.MustCompile(`^\s*(-{2,}|\/\*|\s*)*\s*([A-Z_]+)(.*)$`)

type outputInstruction struct {
	_type  instructionPrefix
	values [][]string
}

func getInstructions(lines []parsedLine) (*outputInstruction, error) {
	instructionPrefixMap := buildMapPrefix()

	instrs := make([]*outputInstruction, 0)
	rowsInstrs := outputInstruction{
		_type: instructionPrefixRow,
	}

	uniquePrefixes := make(map[instructionPrefix]struct{})

	for _, pline := range lines {
		prefixes := rgxInstructionPrefix.FindStringSubmatch(pline.line)

		if len(prefixes) < 4 { //nolint:mnd // 4 is the minimum number of matches
			continue
		}

		prefix := prefixes[2]
		content := prefixes[3]

		prefixType, ok := instructionPrefixMap[prefix]
		if !ok {
			continue
		}

		if _, ok := uniquePrefixes[prefixType]; ok && prefixType.Allowance() == prefixAllowanceSingle {
			return nil, fmt.Errorf("duplicate instruction: %s", prefixType)
		}

		uniquePrefixes[prefixType] = struct{}{}

		switch prefixType {
		case instructionPrefixStartTest, instructionPrefixEndTest:
			continue
		case instructionPrefixCount:
			counts, err := extractCount(content)
			if err != nil {
				return nil, fmt.Errorf("unable to extract count: %w", err)
			}

			instrs = append(instrs, &outputInstruction{
				_type:  prefixType,
				values: counts,
			})

		case instructionPrefixFile:
			rows, err := extractFile(content)
			if err != nil {
				return nil, fmt.Errorf("unable to extract file: %w", err)
			}

			instrs = append(instrs, &outputInstruction{
				_type:  prefixType,
				values: rows,
			})

		case instructionPrefixRow:
			row, err := extractRow(content)
			if err != nil {
				return nil, fmt.Errorf("unable to extract row: %w", err)
			}

			rowsInstrs.values = append(rowsInstrs.values, row)

		default:
			return nil, fmt.Errorf("unknown instruction prefix: %s", prefix)
		}
	}

	if len(rowsInstrs.values) > 0 {
		instrs = append(instrs, &rowsInstrs)
	}

	if err := checkCombinedInstructions(instrs); err != nil {
		return nil, fmt.Errorf("error checking combined instructions: %w", err)
	}

	if len(instrs) == 0 {
		return nil, fmt.Errorf("no instructions found")
	}

	if len(instrs) > 1 {
		return nil, fmt.Errorf("multiple instructions found")
	}

	return instrs[0], nil
}

func checkValidUniquePrefixes(uniquePrefixes map[instructionPrefix]struct{}) error {
	// Can't have any combinations of row, file or count instructions together
	_, foundRow := uniquePrefixes[instructionPrefixRow]
	_, foundFile := uniquePrefixes[instructionPrefixFile]
	_, foundCount := uniquePrefixes[instructionPrefixCount]

	if foundRow && foundFile && foundCount {
		return fmt.Errorf("can't have both ROW, FILE and COUNT instructions")
	}

	if foundRow && foundFile {
		return fmt.Errorf("can't have both ROW and FILE instructions")
	}

	if foundRow && foundCount {
		return fmt.Errorf("can't have both ROW and COUNT instructions")
	}

	if foundFile && foundCount {
		return fmt.Errorf("can't have both FILE and COUNT instructions")
	}

	return nil
}

func extractCount(content string) ([][]string, error) {
	content = strings.TrimSpace(content)

	splitted := strings.Fields(content)

	if len(splitted) == 0 {
		return nil, fmt.Errorf("empty content")
	}

	results := make([][]string, 0, len(splitted))

	for _, s := range splitted {
		results = append(results, []string{s})
	}

	return results, nil
}

func extractRow(content string) ([]string, error) {
	content = strings.TrimSpace(content)

	reader := csv.NewReader(strings.NewReader(content))
	reader.LazyQuotes = true

	// Since it's a single line, we can directly read one record
	record, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV content: %w", err)
	}

	results := make([]string, 0, len(record))
	results = append(results, record...)

	return results, nil
}

func extractFile(content string) ([][]string, error) {
	content = strings.TrimSpace(content)

	file, err := os.Open(filepath.Clean(content))
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close() //nolint:errcheck // we don't care about the error here

	csvReader := csv.NewReader(file)
	csvReader.LazyQuotes = true

	var results [][]string

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break // End of file reached, stop reading
		}

		if err != nil {
			return nil, fmt.Errorf("error reading CSV content: %w", err)
		}

		results = append(results, record)
	}

	return results, nil
}
