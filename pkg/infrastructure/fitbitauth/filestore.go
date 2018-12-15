package fitbitauth

import (
	"io/ioutil"

	"github.com/pkg/errors"

	dfba "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
)

type fileStore struct{}

func NewFileStore() dfba.Store {
	return &fileStore{}
}

func (fs *fileStore) WriteAuthCode(authCode string) error {
	authCodeByte := []byte(authCode)
	err := ioutil.WriteFile("fitbit-auth-code.txt", authCodeByte, 0444)
	if err != nil {
		return errors.Wrap(err, "Error while writing to file: ")
	}
	return nil
}

func (fs *fileStore) ReadAuthCode() (string, error) {
	return "", nil
}
