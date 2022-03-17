package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/KDF5000/nomo/domain/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dsn = "root:root@nomo@tcp(127.0.0.1:3306)/nomo?charset=utf8mb4&parseTime=True&loc=Local"

	userInfo = `{"user_id": "xxx", "open_id": "xxxx", "union_id: "xxx"}`
	pageInfo = `{"notion_secret_key": "xxxx", "notion_page_id": "xxx"}`
)

func TestBindInfoRepo(t *testing.T) {
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		t.Fatal(err)
	}

	repo := NewBindInfoRepo(db)
	db.AutoMigrate(&entity.BindInfo{})
	bindInfo := entity.BindInfo{
		UserPlatform: 1,
		UnionUserID:  "xxxx",
		UserInfo:     userInfo,
		BindPlatform: 1,
		PageInfo:     "{\"page_id\": \"12345678\"}",
	}
	if err := repo.UpdateOrInsert(context.TODO(), &bindInfo); err != nil {
		t.Fatal(err)
	}

	newBindInfo, err := repo.GetBindInfoByUnionUserID(context.TODO(),
		"xxx")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Create: %+v", *newBindInfo)

	time.Sleep(2 * time.Second)
	bindInfo.PageInfo = pageInfo
	if err := repo.UpdateOrInsert(context.TODO(), &bindInfo); err != nil {
		t.Fatal(err)
	}

	newBindInfo, err = repo.GetBindInfoByUnionUserID(context.TODO(),
		"xxx")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Update: %+v", *newBindInfo)
}
