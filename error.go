package main

type NoSessionFound struct{}

func (e *NoSessionFound) Error() string { return "No session found" }
