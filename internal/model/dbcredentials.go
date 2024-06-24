package model

// DBCredentials stores the information we need to connect with a database.
type DBCredentials struct {
	User string
	Pass string
	Name string
	Host string
	Port int
}
