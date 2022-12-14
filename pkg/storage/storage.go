package storage

type RawPost struct {
	Title   string `xml:"title"`       // заголовок публикации
	Link    string `xml:"link"`        // ссылка на источник
	PubTime string `xml:"pubDate"`     // время публикации в формате строки
	Content string `xml:"description"` // содержание публикации
}

type Post struct {
	ID            int    `json:"id" ;sql:"id"`                   // номер записи
	Title         string `json:"title" ;sql:"title"`             // заголовок публикации
	Content       string `json:"content" ;sql:"content"`         // содержание публикации
	PubTime       int64  `json:"pub_time" ;sql:"pub_time"`       // время публикации в формате Unix
	Link          string `json:"link" ;sql:"link"`               // ссылка на источник
	SourceXmlLink string `json:"source_link" ;sql:"source_link"` // Идентификатор источника данных
}

type PostsInterface interface {
	Posts(countPosts int) ([]Post, error)                                    // получение всех новостей
	AddPost(Post) error                                                      // создание новой новости
	RunInsertPosts()                                                         // Добавление новости из канала от RSS
	GetLastPubDateForSources(sourceSlice []string) (map[string]int64, error) // Получение словаря последних данным по источникам RSS
}
