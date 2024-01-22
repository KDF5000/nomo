package notion

import (
	"fmt"
	"strings"
	"time"

	"github.com/KDF5000/nomo/infrastructure/utils"
	"github.com/KDF5000/notion-sdk-go/core"
)

type NotionClient struct{}

func (c *NotionClient) AppendBlock(notionKey, pageId, content string) error {
	client, err := core.NewClient(&core.Option{SecretKey: notionKey})
	if err != nil {
		return err
	}

	if pageId == "" {
		return fmt.Errorf("invalid content")
	}

	page, err := client.RetrivePage(pageId)
	if err != nil {
		return err
	}

	// 2020-08-12T02:12:33.231Z
	lastEditTime, err := time.Parse(time.RFC3339, page.LastEditedTime)
	if err != nil {
		return err
	}

	// ignore error
	children, _, _, _ := client.RetriveBlockChildren(pageId, "", 1)
	// log.Infof("children: %+v", children)
	var blocks []*core.Block
	if lastEditTime.Local().Day() != time.Now().Day() || len(children) == 0 {
		var block core.Block
		block.Object = core.OBJECT_BLOCK
		block.Type = core.BLOCK_HEADING3
		var heading3Block core.HeadingBlobck
		date := time.Now().Format("2006-01-02")
		heading3Block.Text = append(heading3Block.Text, core.RichTextObject{
			Type: core.TYPE_TEXT,
			Text: &core.TextObject{
				Content: date,
			},
		})
		block.Heading3Block = &heading3Block
		blocks = append(blocks, &block)
	}

	var bulletedItem core.ListItemBlock
	bulletedItem.Text = append(bulletedItem.Text,
		core.RichTextObject{
			Type: core.TYPE_TEXT,
			Text: &core.TextObject{
				Content: content,
			},
		})

	blocks = append(blocks, &core.Block{
		Object:                core.OBJECT_BLOCK,
		Type:                  core.BLOCK_BULLETED_LIST_ITEM,
		BulletedListItemBlock: &bulletedItem,
	})

	err = client.AppendBlock(pageId, blocks)
	if err != nil {
		return err
	}

	return nil
}

func (c *NotionClient) getTitle(elements []utils.ContentElement) string {
	if len(elements) == 0 {
		return ""
	}

	var start, end int
	for i, elem := range elements {
		if !elem.IsTag && len(strings.TrimSpace(elem.Text)) != 0 {
			start = i
			break
		}
	}

	for i := len(elements) - 1; i >= 0; i-- {
		if !elements[i].IsTag && len(strings.TrimSpace(elements[i].Text)) != 0 {
			end = i
			break
		}
	}

	var title string
	for i := start; i <= end; i++ {
		title = fmt.Sprintf("%s%s", title, elements[i].Text)
	}
	return strings.Replace(title, "\n", " ", -1)
}

func (c *NotionClient) AddNewPage2Database(notionKey, dbId, content string) error {
	client, err := core.NewClient(&core.Option{SecretKey: notionKey})
	if err != nil {
		return err
	}

	var page core.Page
	page.Parent = core.ParentObject{
		DatabaseID: dbId,
	}

	page.Properties = make(map[string]core.PropertyValue)
	elements := utils.ScanContent(strings.TrimSpace(content))
	var tagObj core.MultiSelectObject
	var contentBlock core.ParagraphBlock
	for _, elem := range elements {
		color := "default"
		if elem.IsTag {
			tagObj = append(tagObj, core.SelectOption{Name: elem.Text[1:]})
			color = "blue"
		}

		obj := core.RichTextObject{
			Type: core.TYPE_TEXT,
			Text: &core.TextObject{
				Content: elem.Text,
			},
			Annotations: &core.AnnotationObject{
				Bold:  true,
				Code:  elem.IsTag,
				Color: color,
			},
		}

		contentBlock.Text = append(contentBlock.Text, obj)
	}

	page.Properties["Name"] = core.PropertyValue{
		Type: core.TYPE_TITLE,
		TitleObject: &core.RichTextArrary{
			core.RichTextObject{
				Type: core.TYPE_TEXT,
				Text: &core.TextObject{
					Content: c.getTitle(elements),
				},
			}},
	}

	var textBlock core.Block
	textBlock.Object = core.OBJECT_BLOCK
	textBlock.Type = core.BLOCK_PARAGRAPH
	textBlock.ParagraphBlock = &contentBlock
	page.Children = append(page.Children, textBlock)

	if len(tagObj) > 0 {
		page.Properties["Tags"] = core.PropertyValue{
			Type:        core.TYPE_MULTI_SELECT,
			MultiSelect: &tagObj,
		}
	}

	return client.CreatePage(&page)
}
