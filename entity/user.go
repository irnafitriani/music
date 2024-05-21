package entity

import "github.com/thedevsaddam/govalidator"

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *User) Rules() govalidator.MapData {
	return govalidator.MapData{
		"email":    []string{"required"},
		"password": []string{"required"},
	}
}
