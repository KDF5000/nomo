package notion

import "testing"

const (
	SecretKey  = "xxxx"
	DatabaseID = "xxx"
)

func TestNotionCreatePage(t *testing.T) {
	cases := []string{
		"#科技 technology change our life!",
		"这是一条没有标签的memo",
		"使用#欢迎 来给内容添加任意标签, 数量不限(注意标签和正文中间应该有个空格哦)",
	}

	client := &NotionClient{}
	for _, content := range cases {
		err := client.AddNewPage2Database(SecretKey, DatabaseID, content)
		if err != nil {
			t.Fatal(err)
		}
	}
}
