package parser

// envError is used to define constant errors.
type envError string

// Error implements the error interface.
func (s envError) Error() string {
	return string(s)
}

const (
	// ErrNoDBHost is returned when the DB_HOST environment variable is not set.
	ErrNoDBHost = envError("DB_HOST is not set")
	// ErrNoDBPort is returned when the DB_PORT environment variable is not set.
	ErrNoDBPort = envError("DB_PORT is not set")
	// ErrDBPortNotNumber is returned when the DB_PORT environment variable is not a number.
	ErrDBPortNotNumber = envError("DB_PORT is not a number")
	// ErrNoDBUser is returned when the DB_USER environment variable is not set.
	ErrNoDBUser = envError("DB_USER is not set")
	// ErrNoDBPass is returned when the DB_PASSWORD environment variable is not set.
	ErrNoDBPass = envError("DB_PASSWORD is not set")
	// ErrNoDBName is returned when the DB_NAME environment variable is not set.
	ErrNoDBName = envError("DB_NAME is not set")
	// ErrNoSQLFile is returned when the SQL_FILE environment variable is not set.
	ErrNoSQLFile = envError("SQL_FILE is not set")
)

type groupError string

func (s groupError) Error() string {
	return string(s)
}

const (
	// ErrInstructionsUnexpectedStart is returned when an instructions group is started inside another instructions group.
	ErrInstructionsUnexpectedStart = groupError("unexpected start of group inside instructions group")
	// ErrsStatementUnexpectedEnd is returned when an instructions group is ended inside a statement group.
	ErrsStatementUnexpectedEnd = groupError("unexpected end of group inside statement group")
)

type runError string

func (s runError) Error() string {
	return string(s)
}

const (
	// ErrParseFile is returned when the file cannot be parsed.
	ErrParseFile = runError("unable to parse file")
	// ErrUnexpectedInstruction is returned when an unexpected instruction is found.
	ErrUnexpectedInstruction = runError("unexpected instruction")
	// ErrUnexpectedStatement is returned when an unexpected statement is found.
	ErrUnexpectedStatement = runError("unexpected statement")
	// ErrUnexpectedGroupType is returned when an unexpected group type is found.
	ErrUnexpectedGroupType = runError("unexpected group type")
)

type sortError string

func (s sortError) Error() string {
	return string(s)
}

const (
	// ErrDifferentRowCount is returned when the actual and expected row counts differ.
	ErrDifferentRowCount = sortError("different row count")
	// ErrDifferentColumnCount is returned when the actual and expected column counts differ.
	ErrDifferentColumnCount = sortError("different column count")
	// ErrAnyNotNullButEmpty is returned when the actual value is empty but the expected value is not.
	ErrAnyNotNullButEmpty = sortError("ANY_NOT_NULL but empty")
)
