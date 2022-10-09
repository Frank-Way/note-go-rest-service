package note

type Note struct {
	Id     uint   `db:"id" json:"id"`
	Title  string `db:"title" json:"title"`
	Text   string `db:"text" json:"text"`
	Author string `db:"author" json:"author"`
}

type Notes = []Note
