package application

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/Jeffail/tunny"
	"github.com/KDF5000/nomo/infrastructure/convertor"
	itemplate "github.com/KDF5000/nomo/interfaces/template"
	"github.com/KDF5000/pkg/log"
)

type IPosterApp interface {
	GenPoster(ctx context.Context, id uint, data interface{}) ([]byte, error)
}

type posterApp struct {
	pool *tunny.Pool
}

func NewPosterApp(workerNum int) *posterApp {
	return &posterApp{
		pool: tunny.NewFunc(workerNum, convertor.ConvertHandler),
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

	cfg := convertor.ConvConfig{
		Ctx:    ctx,
		Url:    server.URL,
		Params: convertor.DefaultHtmlImageParams,
	}

	cfg.Params.Selector = "div.share-nomo"
	out, err := app.pool.ProcessCtx(ctx, cfg)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	dto, _ := out.(convertor.ConvertOutput)
	return dto.Buf, nil
}
