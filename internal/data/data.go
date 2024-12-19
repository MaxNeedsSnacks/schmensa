package data

type Mensa struct {
	Id   string    `json:"id"`
	Name LocString `json:"name"`
}

type Menu []Food

type Food struct {
	Title     LocString `json:"title"`
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
	CounterNames LocString `json:"counterNames"`
}

type LocString struct {
	De string `json:"de"`
	En string `json:"en"`
}
