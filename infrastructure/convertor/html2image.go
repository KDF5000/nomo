package convertor

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

var (
	DefaultViewportWidth         int64 = 390
	DefaultViewportHeight        int64 = 0
	DefaultScale                 int64 = 6
	DefaultFullScreenshotQuality int64 = 60
)

type Html2ImageParams struct {
	ViewportWidth         int64
	ViewportHeight        int64
	Scale                 int64
	Mobile                bool
	FullScreenshotQuality int64
	Selector              string
}

var DefaultHtmlImageParams = Html2ImageParams{
	ViewportWidth:         DefaultViewportWidth,
	ViewportHeight:        DefaultViewportHeight,
	Mobile:                false,
	Scale:                 DefaultScale,
	FullScreenshotQuality: DefaultFullScreenshotQuality,
}

type Html2Image struct {
	conf           ConvConfig
	convertElapsed time.Duration
}

func NewHtml2ImageConvertor(cfg ConvConfig) *Html2Image {
	return &Html2Image{
		conf: cfg,
	}
}

func (c *Html2Image) GetConvertElapsed() time.Duration {
	return c.convertElapsed
}

func (c *Html2Image) Convert(ctx context.Context) ([]byte, error) {
	start := time.Now()
	defer func() {
		c.convertElapsed = time.Since(start)
	}()

	var options []chromedp.ExecAllocatorOption
	options = append(options, chromedp.CombinedOutput(log.Writer()))
	options = append(options, chromedp.DefaultExecAllocatorOptions[:]...)

	options = append(options, chromedp.DisableGPU)
	options = append(options, chromedp.Flag("ignore-certificate-errors", true))
	actX, aCancel := chromedp.NewExecAllocator(ctx, options...)
	defer aCancel()

	ctx, cancel := chromedp.NewContext(actX)
	defer cancel()
	tasks := []chromedp.Action{
		chromedp.Navigate(c.conf.Url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			return emulation.SetDeviceMetricsOverride(c.conf.Params.ViewportWidth,
				c.conf.Params.ViewportHeight,
				float64(c.conf.Params.Scale),
				c.conf.Params.Mobile).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).Do(ctx)
		}),
	}

	var buf []byte
	if c.conf.Params.Selector == "" {
		tasks = append(tasks, chromedp.FullScreenshot(&buf, int(c.conf.Params.FullScreenshotQuality)))
	} else {
		tasks = append(tasks, chromedp.Screenshot(c.conf.Params.Selector, &buf, chromedp.NodeVisible))
	}

	if err := chromedp.Run(ctx, tasks...); err != nil {
		return nil, err
	}

	return buf, nil
}
