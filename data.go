package main

type mensa struct {
	Id   string    `json:"id"`
	Name localized `json:"name"`
}

type menu []food

type food struct {
	Title     localized `json:"title"`
	Type      []string  `json:"type"`
	Additives []string  `json:"additives"`
	Category  string    `json:"category"`
	Price     struct {
		Student string `json:"student"`
		Staff   string `json:"staff"`
		Guest   string `json:"guest"`
	} `json:"price"`
	DispoId      string    `json:"dispoId"`
	Counter      string    `json:"counter"`
	Position     int       `json:"position"`
	CounterNames localized `json:"counterNames"`
}

type localized struct {
	De string `json:"de"`
	En string `json:"en"`
}
