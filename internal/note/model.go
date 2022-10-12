package note

type Note struct {
	Id     uint   `db:"id" json:"id"`
	Title  string `db:"title" json:"title"`
	Text   string `db:"text" json:"text"`
	Author string `db:"author" json:"author"`
}

type Notes = []Note

func NewNote(login string, dto CreateNoteDTO) Note {
	return Note{
		Title:  dto.Title,
		Text:   dto.Text,
		Author: login,
	}
}

func UpdateNote(id uint, login string, dto UpdateNoteDTO) Note {
	return Note{
		Id:     id,
		Title:  dto.Title,
		Text:   dto.Text,
		Author: login,
	}
}
