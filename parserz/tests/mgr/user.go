package mgr

import (
	"context"
	"tests/store"
)

type UserInfo struct {
	Id   int64
	Name string
}

type UserMgr struct {
}

func (u *UserMgr) FindById(ctx context.Context, id int64) *store.User {
	return nil
}

func (u *UserMgr) FindByIdMap(idMap map[int64]UserInfo) []*store.User {
	return nil
}

func (u *UserMgr) FindByIds(idList ...int64) []*store.User {
	return nil
}

func (u *UserMgr) FindByIdList(idList []int64) []*store.User {
	return nil
}

func (u *UserMgr) Save(user *store.User) (int64, error) {
	return 0, nil
}
