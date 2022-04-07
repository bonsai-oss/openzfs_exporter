package main

type arrayFlags []string

func (af *arrayFlags) String() string {
	return ""
}

func (af *arrayFlags) Set(value string) error {
	*af = append(*af, value)
	return nil
}
