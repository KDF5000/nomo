package lark_doc

import (
	"context"
	"time"

	"github.com/KDF5000/pkg/larkdoc"

	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
)

type LarkDocV2Wrapper struct{}

func (c *LarkDocV2Wrapper) findInsertIndex(iterator *larkdocx.ListDocumentBlockIterator) (int, bool, error) {
	hasDate := false
	curDate := time.Now().Format("2006-01-02")
	firstBlock := true

	index := -1
	var retErr error
	for {
		hasMore, block, err := iterator.Next()
		if err != nil {
			retErr = err
			break
		}

		if block == nil {
			break
		}

		// Page title
		if *block.BlockType == larkdoc.DocxBlockTypePage {
			continue
		}

		index++
		if firstBlock {
			if *block.BlockType != larkdoc.DocxBlockTypeHeading3 ||
				len(block.Heading3.Elements) == 0 ||
				*block.Heading3.Elements[0].TextRun.Content != curDate {
				hasDate = false
				break
			}

			// has current date block, need to find last block today
			firstBlock = false
			hasDate = true
			continue
		}

		// find next date block or has no more blocks
		if !hasMore || *block.BlockType == larkdoc.DocxBlockTypeHeading3 {
			break
		}
	}

	return index, hasDate, retErr
}

func (c *LarkDocV2Wrapper) InsertBlock(ctx context.Context, appId, secret string, page, content string) error {
	doc := larkdoc.NewLarkDocV2(larkdoc.DocOption{
		AppID:     appId,
		AppSecret: secret,
	})

	baseInfo, err := doc.GetBasicInfo(ctx, page)
	if err != nil {
		return err
	}

	iter, err := doc.BlocksIterator(ctx, page)
	if err != nil {
		return err
	}

	index, hasDate, err := c.findInsertIndex(iter)
	if err != nil {
		return err
	}

	var blocks []*larkdocx.Block
	if !hasDate {
		dateText := larkdocx.NewTextRunBuilder().Content(time.Now().Format("2006-01-02")).Build()
		text := larkdocx.NewTextBuilder().Elements([]*larkdocx.TextElement{
			larkdocx.NewTextElementBuilder().TextRun(dateText).Build(),
		}).Build()
		dateBlock := larkdocx.NewBlockBuilder()
		dateBlock.BlockType(larkdoc.DocxBlockTypeHeading3)
		dateBlock.Heading3(text)
		blocks = append(blocks, dateBlock.Build())
	}

	memoText := larkdocx.NewTextRunBuilder().Content(content)
	memo := larkdocx.NewTextBuilder().Elements([]*larkdocx.TextElement{
		larkdocx.NewTextElementBuilder().TextRun(memoText.Build()).Build(),
	}).Build()

	memoBlock := larkdocx.NewBlockBuilder()
	memoBlock.BlockType(larkdoc.DocxBlockTypeBullet)
	memoBlock.Bullet(memo)
	blocks = append(blocks, memoBlock.Build())

	return doc.InsertBlock(ctx, page, *baseInfo.RevisionId, index, blocks)
}
