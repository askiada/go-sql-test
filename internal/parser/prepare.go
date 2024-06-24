package parser

func PrepairPair(p pair) (pair, error) {
	sortRows(p.Actual)
	sortRows(p.Expected)

	actual, expected, err := replaceKeywords(p.Actual, p.Expected)
	if err != nil {
		return pair{}, err
	}

	return pair{
		Actual:   actual,
		Expected: expected,
	}, nil
}
