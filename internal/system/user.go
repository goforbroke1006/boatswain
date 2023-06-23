package system

import (
	"os/user"
)

func MustGetCurrentUsername() string {
	currUser, currUserErr := user.Current()
	if currUserErr != nil {
		panic(currUserErr)
	}
	return currUser.Username
}
