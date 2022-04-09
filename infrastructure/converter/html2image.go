package converter

import (
	"context"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

type Html2Image struct {
	buf            []byte
	convertElapsed time.Duration
}

func (ins *Html2Image) GetConvertElapsed() time.Duration {
	return ins.convertElapsed
}

func (ins *Html2Image) Convert(ctx context.Context, url string, sel string) ([]byte, error) {
	start := time.Now()
	defer func() {
		ins.convertElapsed = time.Since(start)
	}()

	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	tasks := []chromedp.Action{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			return emulation.SetDeviceMetricsOverride(390, 844, 6, true).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).Do(ctx)
		}),
	}

	if sel == "" {
		tasks = append(tasks, chromedp.FullScreenshot(&ins.buf, 90))
	} else {
		tasks = append(tasks, chromedp.Screenshot(sel, &ins.buf, chromedp.NodeVisible))
	}

	if err := chromedp.Run(ctx, tasks...); err != nil {
		return nil, err
	}

	return ins.buf, nil
}
