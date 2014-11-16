package main

import (
	"fmt"
)

type User interface {
	SayName() string
}

type UserDefault struct {
	FirstName string
	LastName  string
}

type Customer struct {
	UserDefault
	Kind int
}

func (u *UserDefault) SayName() string {
	str := fmt.Sprintf("My name is %s %s", u.FirstName, u.LastName)
	return str
}

func (c *Customer) SayName() string {
	str := fmt.Sprintf("My name is %s %s, kind is %d.", c.FirstName, c.LastName, c.Kind)
	return str
}

func main() {

	var u User
	// need to set Pointer to UserDefault and Customer 
	u1 := &UserDefault{FirstName: "Taro", LastName: "Yamada"}
	c1 := &Customer{Kind: 1, UserDefault:UserDefault{FirstName: "Jiro", LastName: "Yamada"}}
	
	u = u1
	fmt.Println(u.SayName())
	u = c1
	fmt.Println(u.SayName())

}
