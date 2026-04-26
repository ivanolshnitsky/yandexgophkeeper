package storage

// Credential хранит логин/пароль
type Credential struct {
	ID       string
	User     string
	Login    string
	Password string
	Meta     string
}
