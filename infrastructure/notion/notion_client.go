package notion

import (
	"fmt"
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

func (c *NotionClient) AddNewPage2Database(notionKey, dbId, content string) error {
	client, err := core.NewClient(&core.Option{SecretKey: notionKey})
	if err != nil {
		return err
	}

	var page core.Page
	page.Parent = core.ParentObject{
		DatabaseID: dbId,
	}
	var textBlock core.Block
	textBlock.Object = core.OBJECT_BLOCK
	textBlock.Type = core.BLOCK_PARAGRAPH
	var contentBlock core.ParagraphBlock
	contentBlock.Text = append(contentBlock.Text, core.RichTextObject{
		Type: core.TYPE_TEXT,
		Text: &core.TextObject{
			Content: content,
		},
		Annotations: &core.AnnotationObject{
			Bold: true,
		},
	})
	textBlock.ParagraphBlock = &contentBlock
	page.Children = append(page.Children, textBlock)
	page.Properties = make(map[string]core.PropertyValue)
	page.Properties["Name"] = core.PropertyValue{
		Type:        core.TYPE_TITLE,
		TitleObject: &core.RichTextArrary{},
	}

	tags := utils.RetriveTags(content)
	if len(tags) > 0 {
		var tagObj core.MultiSelectObject
		for _, tag := range tags {
			tagObj = append(tagObj, core.SelectOption{Name: tag})
		}
		page.Properties["Tags"] = core.PropertyValue{
			Type:        core.TYPE_MULTI_SELECT,
			MultiSelect: &tagObj,
		}
	}

	return client.CreatePage(&page)
}
