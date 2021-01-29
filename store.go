package main

type Store struct {
	Identity string  `json:"identity"`
	Entries  Entries `json:"entries"`
}

type Entries []Entry

type Entry struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
