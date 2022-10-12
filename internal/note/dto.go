package note

type CreateNoteDTO struct {
	Title string `db:"title" json:"title"`
	Text  string `db:"text" json:"text"`
}

type UpdateNoteDTO struct {
	Title string `db:"title" json:"title"`
	Text  string `db:"text" json:"text"`
}
