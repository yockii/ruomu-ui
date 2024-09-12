package domain

import "github.com/yockii/ruomu-ui/model"

type Group struct {
	model.MaterialComponentGroup
	Components []Component `json:"components"`
}

type Component struct {
	model.MaterialComponent
	MetaInfo ComponentMetaInfo `json:"metaInfo"`
}

// ComponentMetaInfo components' schema json
type ComponentMetaInfo struct {
	IsContainer bool             `json:"isContainer,omitempty"`
	Props       []PropGroup      `json:"props,omitempty"`
	Events      map[string]Event `json:"events,omitempty"`
	Slots       []string         `json:"slots,omitempty"`
}

type Property struct {
	Label        string            `json:"label,omitempty"`
	Description  string            `json:"description,omitempty"`
	Name         string            `json:"name,omitempty"`
	Type         string            `json:"type,omitempty"`
	Required     bool              `json:"required"`
	DefaultValue interface{}       `json:"defaultValue,omitempty"`
	Widget       Widget            `json:"widget,omitempty"`
	Rules        map[string]string `json:"rules,omitempty"`
}

type WidgetComponent string

const (
	WidgetComponentSwitch   WidgetComponent = "switch"
	WidgetComponentColor    WidgetComponent = "color"
	WidgetComponentSelect   WidgetComponent = "select"
	WidgetComponentInput    WidgetComponent = "input"
	WidgetComponentNumber   WidgetComponent = "number"
	WidgetComponentSlider   WidgetComponent = "slider"
	WidgetComponentFunction WidgetComponent = "function"
)

type Widget struct {
	Component WidgetComponent        `json:"component,omitempty"`
	Inline    bool                   `json:"inline,omitempty"`
	Props     map[string]interface{} `json:"props,omitempty"`
}

type PropGroup struct {
	GroupName  string     `json:"groupName,omitempty"`
	Properties []Property `json:"properties,omitempty"`
}

type Event struct {
	Label       string   `json:"label,omitempty"`
	Description string   `json:"description,omitempty"`
	Name        string   `json:"name,omitempty"`
	Params      []string `json:"params,omitempty"`
}
