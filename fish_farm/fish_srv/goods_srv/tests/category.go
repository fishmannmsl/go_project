package main

import (
	"context"
	"fmt"
	"go_project/fish_farm/fish_srv/goods_srv/global"
	"go_project/fish_farm/fish_srv/goods_srv/model"
	"go_project/fish_farm/fish_srv/goods_srv/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestGetCategorysList(ctx context.Context, req *emptypb.Empty) (*proto.CategoryListResponse, error) {
	var categorys []model.Category
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	for _, category := range categorys {
		fmt.Println(category.Name)
	}
	return nil, nil
}
