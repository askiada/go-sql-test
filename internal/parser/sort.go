package parser

import "sort"

func sortRows(rows [][]string) {
	sort.Slice(rows, func(i, j int) bool {
		// Iterate over each string in the inner slices for comparison
		for k := 0; k < len(rows[i]) && k < len(rows[j]); k++ {
			if rows[i][k] != rows[j][k] {
				// If strings differ, use them to determine order
				return rows[i][k] < rows[j][k]
			}
		}
		// If all compared strings are equal, the shorter slice comes first
		return len(rows[i]) < len(rows[j])
	})
}
