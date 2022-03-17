package utils

import "testing"

func EXPECT_EQ(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestRetriveTags(t *testing.T) {
	cases := []struct {
		Content string
		Tags    []string
	}{
		{
			Content: "没有标签",
			Tags:    []string{},
		},
		{
			Content: "#科技 测试一下",
			Tags:    []string{"科技"},
		},
		{
			Content: "#Tech 测试一下",
			Tags:    []string{"Tech"},
		},
		{
			Content: "#科技无空格标签",
			Tags:    []string{"科技无空格标签"},
		},
		{
			Content: "测试中间#科技 标签",
			Tags:    []string{"科技"},
		},
		{
			Content: "测试结尾标签#科技",
			Tags:    []string{"科技"},
		},
		{
			Content: "#科技 #美食 测试一下两个标签",
			Tags:    []string{"科技", "美食"},
		},
		{
			Content: "# 忘记标签内容",
			Tags:    []string{},
		},
	}

	for _, tc := range cases {
		tags := RetriveTags(tc.Content)
		if !EXPECT_EQ(tags, tc.Tags) {
			t.Fatalf("expected: %+v, got: %+v", tc.Tags, tags)
		}

		t.Logf("content: %s, tags: %+v", tc.Content, tc.Tags)
	}
}
