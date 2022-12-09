package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/ruomu-core/database"
)

type Menu struct {
	Id         int64             `json:"id,omitempty" xorm:"pk"`
	ParentId   int64             `json:"parentId,omitempty" xorm:"comment('父级菜单ID')"`
	Name       string            `json:"name,omitempty" xorm:"comment('菜单名称')"`
	Code       string            `json:"code,omitempty" xorm:"comment('菜单代码')"`
	Icon       string            `json:"icon,omitempty" xorm:"comment('菜单图标')"`
	PageCode   string            `json:"pageCode,omitempty" xorm:"comment('关联的页面代码')"`
	CreateTime database.DateTime `json:"createTime" xorm:"created"`
}

func (_ Menu) TableComment() string {
	return "菜单配置表"
}
func (m *Menu) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	m.Id = j.Get("id").Int()
	m.ParentId = j.Get("parentId").Int()
	m.Name = j.Get("name").String()
	m.Code = j.Get("code").String()
	m.Icon = j.Get("icon").String()
	m.PageCode = j.Get("pageCode").String()

	return nil
}
