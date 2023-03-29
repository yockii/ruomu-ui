package model

import (
	"github.com/tidwall/gjson"

	"github.com/yockii/ruomu-core/database"
)

type Page struct {
	ID            uint64            `json:"id,omitempty,string" xorm:"primaryKey"`
	PageName      string            `json:"pageName,omitempty" xorm:"comment:'页面名称'"`
	PageCode      string            `json:"pageCode,omitempty" xorm:"comment:'页面代码'"`
	ThemeCode     string            `json:"themeCode,omitempty" xorm:"default:'default';comment:'适用主题代码'"` // 默认使用default
	PageRoute     string            `json:"pageRoute,omitempty" xorm:"comment:'页面路由路径'"`
	PageConfig    string            `json:"pageConfig,omitempty" xorm:"type:text;comment:'页面配置'"`
	AuthorizeCode string            `json:"authorizeCode,omitempty" xorm:"comment:'授权代码，不需要权限则留空或设为anon'"`
	CreateTime    database.DateTime `json:"createTime" xorm:"autoCreateTime"`
}

func (p Page) TableComment() string {
	return "页面配置表"
}
func (p *Page) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	p.ID = j.Get("id").Uint()
	p.PageName = j.Get("pageName").String()
	p.PageCode = j.Get("pageCode").String()
	p.ThemeCode = j.Get("themeCode").String()
	p.PageRoute = j.Get("pageRoute").String()
	p.PageConfig = j.Get("pageConfig").String()

	return nil
}
