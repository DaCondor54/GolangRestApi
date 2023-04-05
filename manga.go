package main

type Manga struct {
	BookId  int    `json:"bookId"`
	Title   string `json:"title"`
	Year    int    `json:"year"`
	Volumes int    `json:"volumes"`
}

type Mangas struct {
	MangaArray []Manga
	idCounter  int
}

func (mangas *Mangas) AddManga(manga Manga) {
	manga.BookId = mangas.idCounter
	mangas.idCounter += 1
	mangas.MangaArray = append(mangas.MangaArray, manga)
}

func (mangas Mangas) SameId(id int) int {
	for i, v := range MangaArray.MangaArray {
		if v.BookId == id {
			return i
		}
	}
	return -1
}
