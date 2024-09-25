package domain

import "github.com/yockii/ruomu-ui/model"

type Project struct {
	model.Project
	Store []map[string]any `json:"store"`
	Api   map[string]any   `json:"api"`
}
