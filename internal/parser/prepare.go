package parser

func prepairPair(p pair) (pair, error) {
	sortRows(p.actual)
	sortRows(p.expected)

	actual, expected, err := replaceKeywords(p.actual, p.expected)
	if err != nil {
		return pair{}, err
	}

	return pair{
		actual:   actual,
		expected: expected,
	}, nil
}
