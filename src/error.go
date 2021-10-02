package main

import "fmt"

type EmailInUseError string

func (e EmailInUseError) Error() string {
	return fmt.Sprintf("The email '%v' has already been registered.", string(e))
}
