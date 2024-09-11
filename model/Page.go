package model

import (
	"github.com/tidwall/gjson"
)

type Page struct {
	ID        uint64 `json:"id,omitempty,string" gorm:"primaryKey"`
	ProjectID uint64 `json:"projectId,omitempty,string" gorm:"comment:'项目ID'"`
	Name      string `json:"name,omitempty" gorm:"size:50;comment:'页面名称'"`
	ParentID  uint64 `json:"parentId,omitempty,string" gorm:"comment:'父页面ID'"`
	Route     string `json:"route,omitempty" gorm:"size:50;comment:'页面路由'"`
	Schema    string `json:"schema,omitempty" gorm:"text;comment:'页面配置'"`

	CreateTime int64 `json:"createTime" gorm:"autoCreateTime"`
}

func (p *Page) TableComment() string {
	return "页面配置表"
}

func (p *Page) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	p.ID = j.Get("id").Uint()
	p.ProjectID = j.Get("projectId").Uint()
	p.Name = j.Get("name").String()
	p.ParentID = j.Get("parentId").Uint()
	p.Route = j.Get("route").String()
	p.Schema = j.Get("schema").String()

	return nil
}
