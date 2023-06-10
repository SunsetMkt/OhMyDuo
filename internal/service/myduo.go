// ================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"oh-my-duo/internal/consts"
)

type (
	IMyDuo interface {
		Draw(ctx context.Context, elem consts.MyDuoElements) []byte
	}
)

var (
	localMyDuo IMyDuo
)

func MyDuo() IMyDuo {
	if localMyDuo == nil {
		panic("implement not found for interface IMyDuo, forgot register?")
	}
	return localMyDuo
}

func RegisterMyDuo(i IMyDuo) {
	localMyDuo = i
}
