package model

import (
	"github.com/tidwall/gjson"
)

type Menu struct {
	ID         uint64 `json:"id,omitempty,string" gorm:"primaryKey"`
	ParentID   uint64 `json:"parentId,omitempty,string" gorm:"comment:'父级菜单ID'"`
	Name       string `json:"name,omitempty" gorm:"comment:'菜单名称'"`
	Code       string `json:"code,omitempty" gorm:"comment:'菜单代码'"`
	Icon       string `json:"icon,omitempty" gorm:"comment:'菜单图标'"`
	PageCode   string `json:"pageCode,omitempty" gorm:"comment:'关联的页面代码'"`
	CreateTime int64  `json:"createTime" gorm:"autoCreateTime"`
}

func (_ Menu) TableComment() string {
	return "菜单配置表"
}
func (m *Menu) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	m.ID = j.Get("id").Uint()
	m.ParentID = j.Get("parentId").Uint()
	m.Name = j.Get("name").String()
	m.Code = j.Get("code").String()
	m.Icon = j.Get("icon").String()
	m.PageCode = j.Get("pageCode").String()

	return nil
}
