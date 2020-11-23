package spaceshipmodels

type Spaceship struct {
	Id       int         `orm:"pk; auto" json:"id,omitempty"`
	Name     string      `json:"name,omitempty"`
	Class    string      `json:"class,omitempty"`
	Crew     int         `json:"crew,omitempty"`
	Image    string      `json:"image,omitempty"`
	Value    int         `json:"value,omitempty"`
	Status   string      `json:"status,omitempty"`
	Armament []*Armament `orm:"reverse(many);" json:"armament,omitempty"`
}

type Armament struct {
	Spaceship *Spaceship `orm:"rel(fk)" json:"-"`
	Id        int        `orm:"pk; auto" json:"id,omitempty"`
	Title     string     `json:"title,omitempty"`
	Qty       int        `json:"qty,omitempty"`
}

type Spaceships []Spaceship
