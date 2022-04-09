package application

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/KDF5000/nomo/infrastructure/converter"
	itemplate "github.com/KDF5000/nomo/interfaces/template"
	"github.com/KDF5000/pkg/log"
)

type IPosterApp interface {
	GenPoster(ctx context.Context, id uint, data interface{}) ([]byte, error)
}

type posterApp struct {
	converter *converter.Html2Image
}

func NewPosterApp() *posterApp {
	return &posterApp{
		converter: &converter.Html2Image{},
	}
}

var _ IPosterApp = &posterApp{}

// TODO(kongdefei): add qps limit and refactor to support more template
func (app *posterApp) GenPoster(ctx context.Context, id uint, data interface{}) ([]byte, error) {
	if id != 1 {
		return nil, fmt.Errorf("only support tempalte 1(memo)")
	}

	viewData, ok := data.(itemplate.TPLMemoViewData)
	if !ok {
		return nil, fmt.Errorf("invalid data for the template")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templDir := os.Getenv("TEMPLATE_FILE_PATH")
		tmpl, err := template.ParseFiles(filepath.Join(templDir, "notion_memo.tpl"))
		if err != nil {
			log.Infof("create template failed, err: %s", err)
			return
		}

		tmpl.Execute(w, viewData)
	}))

	return app.converter.Convert(ctx, server.URL, "div.share-nomo")
}
