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
	// manga := Manga{
	// 	mangas.idCounter,
	// 	title,
	// 	year,
	// 	volumes,
	// }
	manga.BookId = mangas.idCounter
	mangas.idCounter += 1
	mangas.MangaArray = append(mangas.MangaArray, manga)
}
