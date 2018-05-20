package utils

type CONFIG_S struct {
	Username string
	Password string
}

var CONFIG CONFIG_S

func init() {
	CONFIG.Username = "lihui02"
	CONFIG.Password = "asdfzxcv"
}
