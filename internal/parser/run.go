package parser

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/askiada/go-sql-test/internal/model"
)

type pair struct {
	Name     string
	Expected [][]string
	Actual   [][]string
}

func Run(ctx context.Context, sqlFile string, db model.DB) ([]pair, error) {
	lines, err := parseFile(sqlFile)
	if err != nil {
		return nil, ErrParseFile
	}

	groups, err := getGroups(lines)
	if err != nil {
		return nil, fmt.Errorf("unable to get groups: %w", err)
	}

	currPair := pair{}

	pairs := []pair{}

	for _, group := range groups {
		switch group._type {
		case groupTypeInstructions:
			instr, err := getInstructions(group.lines)
			if err != nil {
				return nil, fmt.Errorf("unable to get instructions: %w", err)
			}

			currPair.Name = instr.name

			if currPair.Expected == nil {
				currPair.Expected = instr.values
			} else {
				return nil, ErrUnexpectedInstruction
			}

			if currPair.Actual != nil {
				pairs = append(pairs, currPair)
				currPair = pair{}
			}

		case groupTypeStatement:
			rebuildQuery := ""
			for _, line := range group.lines {
				rebuildQuery += line.line + "\n"
			}

			rows, err := db.Query(ctx, rebuildQuery)
			if err != nil {
				return nil, fmt.Errorf("unable to query: %w", err)
			}

			res, err := processRows(rows)
			if err != nil {
				return nil, fmt.Errorf("unable to process rows: %w", err)
			}

			if currPair.Actual == nil {
				currPair.Actual = res
			} else {
				return nil, ErrUnexpectedStatement
			}

			if currPair.Expected != nil {
				pairs = append(pairs, currPair)
				currPair = pair{}
			}
		case groupTypeUnknown:
			return nil, ErrUnexpectedGroupType
		}
	}

	return pairs, nil
}

func processRows(rows pgx.Rows) ([][]string, error) {
	defer rows.Close()

	res := [][]string{}

	for rows.Next() {
		rowAny, err := rows.Values()
		if err != nil {
			return nil, fmt.Errorf("unable to get values: %w", err)
		}

		row := make([]string, 0, len(rowAny))

		for _, v := range rowAny {
			row = append(row, fmt.Sprintf("%v", v))
		}

		res = append(res, row)
	}

	return res, nil
}
