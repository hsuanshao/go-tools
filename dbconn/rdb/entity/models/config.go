package rdbmodel

type Connect struct {
	Driver   string `json:"driver"`
	Network  string `json:"network"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"database_name"`
}
