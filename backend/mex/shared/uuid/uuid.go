package uuid

import uuidBase "github.com/satori/go.uuid"

func MustNewV4() string {
	u, err := uuidBase.NewV4()
	if err != nil {
		panic(err)
	}

	return u.String()
}
