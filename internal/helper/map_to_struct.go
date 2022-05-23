package helper

import "encoding/json"

func MapToStruct(mapData interface{}, stuctData interface{}) error {
	data, err := json.Marshal(mapData)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, stuctData)
	if err != nil {
		return err
	}
	return nil
}
