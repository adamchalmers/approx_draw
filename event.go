package main

type Event struct {
	Name       string            `json:"type"`
	Properties map[string]string `json:"properties",omitempty`
}
