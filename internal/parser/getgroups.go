package parser

type groupType int

const (
	groupTypeUnknown groupType = iota
	groupTypeInstructions
	groupTypeStatement
)

func (gt groupType) String() string {
	switch gt {
	case groupTypeInstructions:
		return "Instructions"
	case groupTypeStatement:
		return "Statement"
	default:
		return "Unknown"
	}
}

type groupLines struct {
	lines []parsedLine
	_type groupType
}

func getGroups(lines []parsedLine) ([]*groupLines, error) {
	var (
		results       []*groupLines
		group         []parsedLine
		currGroupType = groupTypeUnknown
		nextGroupType = groupTypeUnknown
	)

	for _, line := range lines {
		switch currGroupType {
		case groupTypeUnknown:
			switch line.lineType {
			case lineTypeUnknown:
				group = append(group, line)
				nextGroupType = groupTypeStatement
			case lineTypeStartTest:
				group = append(group, line)
				nextGroupType = groupTypeInstructions
			case lineTypeEndTest:
				group = append(group, line)
				gl := &groupLines{
					lines: group,
					_type: currGroupType,
				}
				results = append(results, gl)
				group = []parsedLine{}
				nextGroupType = groupTypeUnknown
			case lineTypeComment:
				group = append(group, line)
			}

		case groupTypeInstructions:
			switch line.lineType {
			case lineTypeUnknown:
				group = append(group, line)
			case lineTypeStartTest:
				return nil, ErrInstructionsUnexpectedStart
			case lineTypeEndTest:
				group = append(group, line)
				gl := &groupLines{
					lines: group,
					_type: currGroupType,
				}
				results = append(results, gl)
				group = []parsedLine{}
				nextGroupType = groupTypeUnknown
			case lineTypeComment:
				group = append(group, line)
			}
		case groupTypeStatement:
			switch line.lineType {
			case lineTypeUnknown:
				group = append(group, line)
			case lineTypeStartTest:
				gl := &groupLines{
					lines: group,
					_type: currGroupType,
				}
				results = append(results, gl)
				group = []parsedLine{line}
				nextGroupType = groupTypeInstructions
			case lineTypeEndTest:
				return nil, ErrsStatementUnexpectedEnd
			case lineTypeComment:
				group = append(group, line)
			}
		}

		currGroupType = nextGroupType
	}

	if len(group) > 0 {
		gl := &groupLines{
			lines: group,
			_type: currGroupType,
		}
		results = append(results, gl)
	}

	return results, nil
}
