package resources

type PostTagResource struct {
	Id      int64    `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}
