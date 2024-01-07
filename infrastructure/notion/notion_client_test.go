package notion

import (
	"os"
	"testing"

	"github.com/KDF5000/nomo/infrastructure/utils"
)

var (
	SecretKey  = os.Getenv("secret_key")
	DatabaseID = os.Getenv("database_id")
)

func TestNotionCreatePage(t *testing.T) {
	cases := []string{
		"#科技 technology change our life!",
		"这是一条没有标签的memo",
		"使用#欢迎 来给内容添加任意标签, 数量不限(注意标签和正文中间应该有个空格哦)",
		"#科技 只是一条科技#美食 memo",
		"有些人喜欢在中间加#标签 然后",
		"这是一个很长很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很很长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长长的内容",
		"#科技 #商业 有两个tag",
		"#科技 #商业 有两个tag和一个尾部tag#尾部",
		"#科技  #商业 有两个tag和一个尾部tag#尾部  ",
		"#科技 #商业 有两个tag和一个尾部tag#空格   #尾部  ",
		"#科技#商业 有两个tag和一个尾部tag#空格   #尾部  ",
	}

	client := &NotionClient{}
	for _, content := range cases {
		err := client.AddNewPage2Database(SecretKey, DatabaseID, content)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestScanContent(t *testing.T) {
	t.Logf("%+v", utils.ScanContent("有些人喜欢在中间加#标签 然后"))
	t.Logf("%v", utils.ScanContent("#科技 只是一条科技#美食 哈哈"))
}
