package model

import (
	"github.com/tidwall/gjson"
)

type Project struct {
	ID          uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	Name        string `json:"name,omitempty" gorm:"size:50;comment:'项目名称'"`
	Description string `json:"description,omitempty" gorm:"size:255;comment:'项目描述'"`
	HomePageID  uint64 `json:"homePageID,omitempty,string" gorm:"comment:'首页ID'"`
	Status      int    `json:"status,omitempty" gorm:"comment:'项目状态'"`
	CreateTime  int64  `json:"createTime" gorm:"autoCreateTime"`
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

type ProjectFrontend struct {
	ID        uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	StoreJson string `json:"-" gorm:"text;comment:'项目级别的变量json'"`
	JsJson    string `json:"-" gorm:"text;comment:'js代码'"`
	CssJson   string `json:"-" gorm:"text;comment:'css代码'"`
	ApiJson   string `json:"-" gorm:"text;comment:'api代码'"`
}

func (p *ProjectFrontend) TableComment() string {
	return "项目前端代码表"
}

func (p *ProjectFrontend) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	p.ID = j.Get("id").Uint()
	p.StoreJson = j.Get("store").String()
	p.JsJson = j.Get("js").String()
	p.CssJson = j.Get("css").String()
	p.ApiJson = j.Get("api").String()

	return nil
}

type ProjectMaterialLibVersion struct {
	ID           uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ProjectID    uint64 `json:"projectId,omitempty,string" gorm:"index;comment:'项目ID'"`
	LibID        uint64 `json:"libId,omitempty,string" gorm:"index;comment:'物料库ID'"`
	LibVersionID uint64 `json:"libVersionId,omitempty,string" gorm:"index;comment:'物料库版本ID'"`
	CreateTime   int64  `json:"createTime" gorm:"autoCreateTime"`
}

func (p *ProjectMaterialLibVersion) TableComment() string {
	return "项目物料库版本表"
}

func (p *ProjectMaterialLibVersion) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	p.ID = j.Get("id").Uint()
	p.ProjectID = j.Get("projectId").Uint()
	p.LibID = j.Get("libId").Uint()
	p.LibVersionID = j.Get("libVersionId").Uint()

	return nil
}
