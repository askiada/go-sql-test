package parser

type Keyword string

const (
	KeywordAny        Keyword = "<K_ANY>"
	KeywordAnyNotNull Keyword = "<K_ANY_NOT_NULL>"
)

func replaceKeywords(actual, expected [][]string) ([][]string, [][]string, error) {
	if len(actual) != len(expected) {
		return nil, nil, ErrDifferentRowCount
	}

	for i := range actual {
		if len(actual[i]) != len(expected[i]) {
			return nil, nil, ErrDifferentColumnCount
		}

		for j := range actual[i] {
			if expected[i][j] == string(KeywordAny) {
				actual[i][j] = string(KeywordAny)
			} else if expected[i][j] == string(KeywordAnyNotNull) {
				if actual[i][j] == "" {
					return nil, nil, ErrAnyNotNullButEmpty
				}

				if actual[i][j] == "null" {
					return nil, nil, ErrAnyNotNullButEmpty
				}

				if actual[i][j] == "NULL" {
					return nil, nil, ErrAnyNotNullButEmpty
				}

				if actual[i][j] == "<nil>" {
					return nil, nil, ErrAnyNotNullButEmpty
				}

				actual[i][j] = string(KeywordAnyNotNull)
			}
		}
	}

	return actual, expected, nil
}
