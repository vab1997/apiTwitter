package models

type SecretManager struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	JWTSign  string `json:"jwtsign"`
	DataBase string `json:"database"`
}
