package main

import (
	"fmt"
)

type ErrorDTO struct {
	Code    string `json:"err_code"`
	Message string `json:"err_msg"`
}

type EmailInUseError string

func (e EmailInUseError) Error() string {
	return fmt.Sprintf("The email '%v' has already been registered.", string(e))
}
