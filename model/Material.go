package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/ruomu-core/database"
)

type MaterialLib struct {
	ID              uint64            `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	Code            string            `json:"code,omitempty" gorm:"size:50;comment:'物料库代码'"`
	Name            string            `json:"name,omitempty" gorm:"size:50;comment:'物料库名称'"`
	PackageName     string            `json:"packageName,omitempty" gorm:"size:50;comment:'npm包名'"`
	Website         string            `json:"website,omitempty" gorm:"size:255;comment:'官网地址'"`
	Description     string            `json:"description,omitempty" gorm:"size:255;comment:'物料库描述'"`
	ThumbnailUrl    string            `json:"thumbnailUrl,omitempty" gorm:"size:255;comment:'缩略图地址'"`
	ActiveVersionID uint64            `json:"activeVersionId,omitempty,string" gorm:"comment:'当前版本ID'"`
	CreateTime      database.DateTime `json:"createTime" gorm:"autoCreateTime"`
}

func (m *MaterialLib) TableComment() string {
	return "物料库表"
}

func (m *MaterialLib) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	m.ID = j.Get("id").Uint()
	m.Code = j.Get("code").String()
	m.Name = j.Get("name").String()
	m.PackageName = j.Get("packageName").String()
	m.Website = j.Get("website").String()
	m.Description = j.Get("description").String()
	m.ThumbnailUrl = j.Get("thumbnailUrl").String()
	m.ActiveVersionID = j.Get("activeVersionId").Uint()

	return nil
}

type MaterialLibVersion struct {
	ID            uint64            `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	LibID         uint64            `json:"libId,omitempty,string" gorm:"comment:'物料库ID'"`
	Version       string            `json:"version,omitempty" gorm:"size:50;comment:'版本号'"`
	PluginUseName string            `json:"pluginUseName,omitempty" gorm:"size:50;comment:'插件使用名称'"`
	CdnJsUrl      string            `json:"cdnJsUrl,omitempty" gorm:"size:255;comment:'CDN JS地址'"`
	CdnCssUrl     string            `json:"cdnCssUrl,omitempty" gorm:"size:255;comment:'CDN CSS地址'"`
	CreateTime    database.DateTime `json:"createTime" gorm:"autoCreateTime"`
}

func (m *MaterialLibVersion) TableComment() string {
	return "物料库版本表"
}

func (m *MaterialLibVersion) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	m.ID = j.Get("id").Uint()
	m.LibID = j.Get("libId").Uint()
	m.Version = j.Get("version").String()
	m.PluginUseName = j.Get("pluginUseName").String()
	m.CdnJsUrl = j.Get("cdnJsUrl").String()
	m.CdnCssUrl = j.Get("cdnCssUrl").String()

	return nil
}

type MaterialComponentGroup struct {
	ID          uint64            `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	LibID       uint64            `json:"libId,omitempty,string" gorm:"comment:'物料库ID'"`
	LibCode     string            `json:"libCode,omitempty" gorm:"size:50;comment:'物料库代码'"`
	Name        string            `json:"name,omitempty" gorm:"size:50;comment:'分组名称'"`
	Description string            `json:"description,omitempty" gorm:"size:255;comment:'分组描述'"`
	OrderNum    int               `json:"orderNum,omitempty" gorm:"comment:'排序'"`
	CreateTime  database.DateTime `json:"createTime" gorm:"autoCreateTime"`
}

func (m *MaterialComponentGroup) TableComment() string {
	return "物料库分组表"
}

func (m *MaterialComponentGroup) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	m.ID = j.Get("id").Uint()
	m.LibID = j.Get("libId").Uint()
	m.LibCode = j.Get("libCode").String()
	m.Name = j.Get("name").String()
	m.Description = j.Get("description").String()
	m.OrderNum = int(j.Get("orderNum").Int())

	return nil
}

type MaterialComponent struct {
	ID          uint64            `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	LibID       uint64            `json:"libId,omitempty,string" gorm:"comment:'物料库ID'"`
	GroupID     uint64            `json:"groupId,omitempty,string" gorm:"comment:'分组ID'"`
	Name        string            `json:"name,omitempty" gorm:"size:50;comment:'组件名称'"`
	Description string            `json:"description,omitempty" gorm:"size:255;comment:'组件描述'"`
	TagName     string            `json:"tagName,omitempty" gorm:"size:50;comment:'组件标签名'"`
	Thumbnail   string            `json:"thumbnail,omitempty" gorm:"size:255;comment:'缩略图'"`
	Schema      string            `json:"schema,omitempty" gorm:"text;comment:'配置json'"`
	CreateTime  database.DateTime `json:"createTime" gorm:"autoCreateTime"`
}

func (m *MaterialComponent) TableComment() string {
	return "物料库组件表"
}

func (m *MaterialComponent) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	m.ID = j.Get("id").Uint()
	m.LibID = j.Get("libId").Uint()
	m.GroupID = j.Get("groupId").Uint()
	m.Name = j.Get("name").String()
	m.Description = j.Get("description").String()
	m.TagName = j.Get("tagName").String()
	m.Thumbnail = j.Get("thumbnail").String()
	m.Schema = j.Get("schema").String()

	return nil
}

type MaterialComponentVersion struct {
	ID           uint64 `json:"id,omitempty,string" gorm:"primaryKey;autoIncrement:false"`
	ComponentID  uint64 `json:"componentId,omitempty,string" gorm:"comment:'组件ID'"`
	LibVersionID uint64 `json:"libVersionId,omitempty,string" gorm:"comment:'物料库版本ID'"`
}

func (m *MaterialComponentVersion) TableComment() string {
	return "物料库组件版本表"
}

func (m *MaterialComponentVersion) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	m.ID = j.Get("id").Uint()
	m.ComponentID = j.Get("componentId").Uint()
	m.LibVersionID = j.Get("libVersionId").Uint()

	return nil
}
