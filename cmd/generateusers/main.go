package main

import (
	"github.com/larwef/lcoin/user"
)

func main() {
	user1 := user.NewUser("user1")
	user2 := user.NewUser("user2")
	user3 := user.NewUser("user3")

	user1.PersistUser()
	user2.PersistUser()
	user3.PersistUser()

	// Just to test that loading works
	user1 = user.LoadUser("user1")
	user2 = user.LoadUser("user2")
	user3 = user.LoadUser("user3")
}
