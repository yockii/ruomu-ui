package model

import (
	"github.com/tidwall/gjson"

	"github.com/yockii/ruomu-core/database"
)

type Page struct {
	Id         int64             `json:"id,omitempty" xorm:"pk"`
	PageName   string            `json:"pageName,omitempty" xorm:"comment('页面名称')"`
	PageCode   string            `json:"pageCode,omitempty" xorm:"comment('页面代码')"`
	ThemeCode  string            `json:"themeCode,omitempty" xorm:"default('default') comment('适用主题代码')"` // 默认使用default
	PageRoute  string            `json:"pageRoute,omitempty" xorm:"comment('页面路由路径')"`
	PageConfig string            `json:"pageConfig,omitempty" xorm:"text comment('页面配置')"`
	CreateTime database.DateTime `json:"createTime" xorm:"created"`
}

func (p Page) TableComment() string {
	return "页面配置表"
}
func (p *Page) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	p.Id = j.Get("id").Int()
	p.PageName = j.Get("pageName").String()
	p.PageCode = j.Get("pageCode").String()
	p.ThemeCode = j.Get("themeCode").String()
	p.PageRoute = j.Get("pageRoute").String()
	p.PageConfig = j.Get("pageConfig").String()

	return nil
}
