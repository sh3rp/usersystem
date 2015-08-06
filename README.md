usersystem written in go

example code

package main

import (
	"fmt"
	"github.com/sh3rp/usersystem"
)

func main() {
	ds := usersystem.New()
	ds.Open()
	user := ds.NewUser("shep")
	user, err := ds.GetUser("shep")
	fmt.Println(user)
	user.Email = "skendall@gmail.com"
	ds.UpdateUser(user)
	user, err = ds.GetUser("shep")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(user.String())
}
