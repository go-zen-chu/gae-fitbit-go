package fitbitauth

type Store interface {
	WriteAuthCode(authCode string) error
	ReadAuthCode() (string, error)
}
