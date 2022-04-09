package converter

import (
	"context"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	c := &Html2Image{}

	buf, err := c.Convert(context.TODO(), "http://baidu.com", "")
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("write file cost %d", c.GetConvertElapsed())
	if err := ioutil.WriteFile("screenshot1.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}
}

func TestSelector(t *testing.T) {
	type Element struct {
		IsTag   bool
		Content string
	}

	type ViewData struct {
		CreatedAt       string
		ContentElements []Element
		UserName        string
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./template.tpl")
		if err != nil {
			log.Printf("create template failed, err: %s", err)
			return
		}

		Elements := []Element{
			{
				IsTag:   true,
				Content: "#思考",
			},
			{
				IsTag:   false,
				Content: "没有标签的flomo还有存在的意义吗？为什么不用notion文档替代呢？这是一个中间的",
			},
			{
				IsTag:   true,
				Content: "#测试",
			},
			{
				IsTag:   false,
				Content: "标签",
			},
		}

		var data ViewData
		data.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
		data.UserName = "DF.K"
		data.ContentElements = Elements

		tmpl.Execute(w, data)
	}))

	c := &Html2Image{}
	buf, err := c.Convert(context.TODO(), server.URL, "div.share-nomo")
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("convert %s succ! write file cost %v", server.URL, c.GetConvertElapsed())
	if err := ioutil.WriteFile("screenshot2.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}
}
