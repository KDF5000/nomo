package lark_doc

import (
	"context"
	"time"

	"github.com/KDF5000/pkg/larkdoc"
)

const (
	DateHeadingLevel = 3
)

type LarkDocWrapper struct{}

func (c *LarkDocWrapper) findFirstIndex(content *larkdoc.DocContent) (uint64, bool, error) {
	firstIndex := content.Title.Location.EndIndex
	hasDate := false
	curDate := time.Now().Format("2006-01-02")
	firstBlock := true

	for i := range content.Body.Blocks {
		block := &content.Body.Blocks[i]
		if block.Type != larkdoc.BLOCK_TYPE_PARAGRAPH {
			continue
		}

		p := block.Paragraph
		if firstBlock {
			if p.Style.HeadingLevel != DateHeadingLevel ||
				p.Elements[0].TextRun.Text != curDate {
				firstIndex = p.Location.StartIndex
				hasDate = false
				break
			}

			firstBlock = false
			hasDate = true
			continue
		}

		if p.Style.HeadingLevel == DateHeadingLevel {
			firstIndex = p.Location.StartIndex
			break
		}

		firstIndex = block.Paragraph.Location.EndIndex
	}

	return firstIndex, hasDate, nil
}

func (c *LarkDocWrapper) InsertBlock(appId, secret string, page, content string) error {
	doc := larkdoc.NewLarkDoc(larkdoc.DocOption{
		AppID:     appId,
		AppSecret: secret,
	})

	oldContent, err := doc.GetContent(context.Background(), page)
	if err != nil {
		return err
	}

	firstIndex, hasDate, err := c.findFirstIndex(oldContent)
	if err != nil {
		return err
	}

	var blocks []larkdoc.Block
	if !hasDate {
		var dateBlock larkdoc.Block
		dateBlock.Type = larkdoc.BLOCK_TYPE_PARAGRAPH
		var paragraph larkdoc.ParagraphBlock
		paragraph.Style.HeadingLevel = DateHeadingLevel
		paragraph.Elements = append(paragraph.Elements,
			larkdoc.ParagraphElement{Type: larkdoc.PARAGRAPH_ELEMENT_TYPE_TEXTRUN,
				TextRun: &larkdoc.TextRun{Text: time.Now().Format("2006-01-02")}})
		dateBlock.Paragraph = &paragraph
		blocks = append(blocks, dateBlock)
	}

	var memo larkdoc.Block
	memo.Type = larkdoc.BLOCK_TYPE_PARAGRAPH
	memo.Paragraph = &larkdoc.ParagraphBlock{}
	memo.Paragraph.Style.List.Type = larkdoc.LIST_TYPE_BULLET
	memo.Paragraph.Style.List.IndentLevel = 1
	memo.Paragraph.Elements = append(memo.Paragraph.Elements,
		larkdoc.ParagraphElement{Type: larkdoc.PARAGRAPH_ELEMENT_TYPE_TEXTRUN,
			TextRun: &larkdoc.TextRun{Text: content}})

	blocks = append(blocks, memo)
	return doc.InsertBlock(context.Background(), page, oldContent.Revision,
		firstIndex, blocks...)
}
