package main

import (
	"Salehaskarzadeh/internal/storee"
	"fmt"
)

func main() {
	var p storee.Password
	err := p.Set("test123")
	if err != nil {
		fmt.Println("Hash error:", err)
	} else {
		fmt.Println("Password hashed successfully")
	}
}
