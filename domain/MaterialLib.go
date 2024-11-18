package domain

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"github.com/yockii/ruomu-ui/model"
)

type MaterialLibCondition struct {
	model.MaterialLib
	ProjectID uint64 `json:"projectId,omitempty,string"`
}

func (c *MaterialLibCondition) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &c.MaterialLib)
	if err != nil {
		return err
	}
	c.ProjectID = gjson.GetBytes(data, "projectId").Uint()
	return nil
}
