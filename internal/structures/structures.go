package structures

type Task struct {
	ID      string  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

var Secret []byte = []byte("!hawk?2021?avangard!")


type Password struct {
	Password string `json:"password"`
}