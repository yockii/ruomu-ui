package model

type Menu struct {
	Id       int64  `json:"id,omitempty" xorm:"pk"`
	ParentId int64  `json:"parentId,omitempty" xorm:"comment('父级菜单ID')"`
	Name     string `json:"name,omitempty" xorm:"comment('菜单名称')"`
	Code     string `json:"code,omitempty" xorm:"comment('菜单代码')"`
	Icon     string `json:"icon,omitempty" xorm:"comment('菜单图标')"`
	PageCode string `json:"pageCode,omitempty" xorm:"comment('关联的页面代码')"`
}

func (_ Menu) TableComment() string {
	return "菜单配置表"
}
