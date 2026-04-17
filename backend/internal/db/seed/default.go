package seed

import (
	"context"
	"fmt"
	"time"

	"blog/internal/ent"
)

func RunDefault(ctx context.Context, client *ent.Client) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}

	if _, err := tx.Comment.Delete().Exec(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err := tx.Article.Delete().Exec(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err := tx.Category.Delete().Exec(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}

	categoryNames := []string{"Go", "Docker", "Cloudflare", "AWS", "Raspberry Pi"}
	categoryIDsByName := make(map[string]uint32, len(categoryNames))
	for _, name := range categoryNames {
		category, err := tx.Category.Create().
			SetName(name).
			Save(ctx)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		categoryIDsByName[name] = category.ID
	}

	publishStartAt := time.Now().Add(-24 * time.Hour)
	categoryNamesByArticleIndex := map[int][]string{
		1: {"Go", "Docker", "Cloudflare", "AWS", "Raspberry Pi"},
		2: {"Cloudflare", "AWS"},
		3: {"Go", "Raspberry Pi"},
	}

	for i := 1; i <= 100; i++ {
		article, err := tx.Article.Create().
			SetTitle(fmt.Sprintf("かさまし記事%d", i)).
			SetBody("記事の件数を増やすためのかさまし記事です。内容は適当です。").
			SetIsPublished(i%2 == 0).
			SetPublishStartAt(publishStartAt).
			Save(ctx)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		categoryNames, ok := categoryNamesByArticleIndex[i]
		if !ok {
			continue
		}

		categoryIDs := make([]uint32, 0, len(categoryNames))
		for _, categoryName := range categoryNames {
			categoryID, ok := categoryIDsByName[categoryName]
			if !ok {
				_ = tx.Rollback()
				return fmt.Errorf("カテゴリが見つかりません: %s", categoryName)
			}
			categoryIDs = append(categoryIDs, categoryID)
		}

		if err := tx.Article.UpdateOneID(article.ID).AddCategoryIDs(categoryIDs...).Exec(ctx); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
