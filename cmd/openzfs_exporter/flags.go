package main

import "fmt"

type arrayFlags []string

func (af *arrayFlags) String() string {
	return fmt.Sprintf("%+v", af)
}

func (af *arrayFlags) Set(value string) error {
	*af = append(*af, value)
	return nil
}
