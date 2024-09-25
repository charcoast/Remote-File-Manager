package main

type Feature struct {
	Prefixes []string
}
type FeatureRegister struct {
	Port     int      `json:"port"`
	Prefixes []string `json:"prefixes"`
}
