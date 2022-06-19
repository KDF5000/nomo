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
	"github.com/KDF5000/nomo/interfaces/proto"
	itemplate "github.com/KDF5000/nomo/interfaces/template"
	"github.com/KDF5000/pkg/log"
)

type IPosterApp interface {
	GenPoster(ctx context.Context, id uint, data interface{}) ([]byte, error)
	Screnshot(ctx context.Context, req *proto.ScreenshotRequst) ([]byte, error)
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
		log.Errorf("failed to process genposter: %v", err)
		return nil, err
	}
	dto, _ := out.(convertor.ConvertOutput)
	return dto.Buf, nil
}

func (app *posterApp) Screnshot(ctx context.Context, req *proto.ScreenshotRequst) ([]byte, error) {
	cfg := convertor.ConvConfig{
		Ctx:    ctx,
		Url:    req.Url,
		Params: convertor.DefaultHtmlImageParams,
	}

	if req.Mobile > 0 {
		cfg.Params.Mobile = true
	}

	if req.Width > 0 {
		cfg.Params.ViewportWidth = req.Width
	}

	if req.Height > 0 {
		cfg.Params.ViewportHeight = req.Height
	}

	if req.Quality > 0 {
		cfg.Params.FullScreenshotQuality = req.Quality
	}

	out, err := app.pool.ProcessCtx(ctx, cfg)
	if err != nil {
		log.Errorf("failed to screenshot, err=%v", err)
	}

	dto, _ := out.(convertor.ConvertOutput)
	return dto.Buf, nil
}
