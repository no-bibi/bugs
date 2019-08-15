package xorm

import (
	"github.com/no-bibi/bugs/fun"
)

type PkString struct {
	Id string `xorm:"_id pk varchar(32) notnull" json:"_id" form:"id"`
}

//唯一pk处理
func (pk *PkString) BeforeInsert() {
	if pk.Id == `` {
		pk.Id, _ = fun.Unique()
	}
}

// demo
//type Account struct {
//	PkString `xorm:"extends"`
//	Itime   int64 `xorm:"created"`
//	Utime   int64 `xorm:"updated"`
//}
