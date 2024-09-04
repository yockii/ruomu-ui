package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/ruomu-core/database"
)

type Project struct {
	ID          uint64            `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	Name        string            `json:"name,omitempty" gorm:"size:50;comment:'项目名称'"`
	Description string            `json:"description,omitempty" gorm:"size:255;comment:'项目描述'"`
	HomePageID  uint64            `json:"homePageID,omitempty,string" gorm:"comment:'首页ID'"`
	Status      int               `json:"status,omitempty" gorm:"comment:'项目状态'"`
	CreateTime  database.DateTime `json:"createTime" gorm:"autoCreateTime"`
}

func (p *Project) TableComment() string {
	return "项目表"
}

func (p *Project) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	p.ID = j.Get("id").Uint()
	p.Name = j.Get("name").String()
	p.Description = j.Get("description").String()
	p.HomePageID = j.Get("homePageID").Uint()
	p.Status = int(j.Get("status").Int())

	return nil
}
