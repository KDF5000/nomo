package proto

type PosterRequest struct {
	UserName  string `form:"user_name" json:"user_name"`
	CreatedAt string `form:"created_at" json:"created_at"`
	Content   string `form:"content" json:"content"`
}

func (req *PosterRequest) IsValid() bool {
	return req.UserName != "" && req.Content != "" && req.CreatedAt != ""
}
