package structures

// swagger:model
type Task struct {
	ID      string  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

var Secret []byte = []byte("!hawk?2021?avangard!")


// swagger:model
type Password struct {
	Password string `json:"password"`
}



// swagger:model
type StatusBadRequest struct {
	Error string `json:"error" example:"invalid request"`
}

// swagger:model
type StatusOK struct {
	Status string `json:"status" example:"ok"`
}

// swagger:model
type StatusInternalServerError struct {
	Error string `json:"error" example:"internal server error"`
}

// swagger:model
type StatusUnauthorized struct {
	Error string `json:"error" example:"unauthorized"`
}

//swagger:model
type StatusNotFound struct {
	Message string `json:"message" example:"Error: not found"`
  }