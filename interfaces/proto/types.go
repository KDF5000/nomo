package proto

import (
	"strings"

	"github.com/KDF5000/pkg/log"
)

type PosterRequest struct {
	UserName  string `form:"user_name" json:"user_name"`
	CreatedAt string `form:"created_at" json:"created_at"`
	Content   string `form:"content" json:"content"`
}

func (req *PosterRequest) IsValid() bool {
	return req.UserName != "" && req.Content != "" && req.CreatedAt != ""
}

type ScreenshotRequst struct {
	Url     string `json:"url" form:"url"`
	Width   int64  `json:"width" form:"width"`
	Height  int64  `json:"height" form:"height"`
	Quality int64  `json:"quality" form:"quality"`
	Mobile  uint8  `json:"mobile" form:"mobile"`
}

func (r *ScreenshotRequst) IsValidUrl() bool {
	if r.Url == "" || !(strings.HasPrefix(r.Url, "http://") ||
		strings.HasPrefix(r.Url, "https://")) {
		log.Infof("url %v, %v, %v", r.Url,
			strings.HasPrefix(r.Url, "http://"), strings.HasPrefix(r.Url, "https://"))
		return false
	}

	return true
}
