package resources

type PostTagResource struct {
	Id      int64    `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type PostByTagResource struct {
	Tag         string        `json:"tag"`
	DetailPosts []DetailPosts `json:"detail_posts"`
}

type DetailPosts struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
