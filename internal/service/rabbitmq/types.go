package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"log"
	"strings"
)

type Receiver func(body []byte)

//COMMON
type Event struct {
	ID       string                 `json:"event_id"`
	Version  int                    `json:"event_version"`
	Name     string                 `json:"event_name"`
	Time     string                 `json:"event_time"`
	Producer string                 `json:"producer"`
	Data     map[string]interface{} `json:"data"`
}

func (e *Event) Validate(eventType string, version int) bool {
	jsonData, _ := json.Marshal(e)
	schema := e.getSchemaPath(eventType, version)
	return e.isValid(jsonData, schema)
}

func (e *Event) getSchemaPath(eventType string, version int) string {
	eName := strings.Replace(eventType, ".", "/", 1)
	return fmt.Sprintf("../../internal/event/schemas/%s/%o.json", eName, version)
	//return fmt.Sprintf("%s/%s/%o.json", BASE_URL, eName, version)
}

func (e *Event) isValid(data []byte, schemaPath string) bool {
	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft4
	sch, err := compiler.Compile(schemaPath)
	if err != nil {
		log.Printf("%#v", err.Error())
		return false
	}

	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		log.Printf("%#v", err.Error())
		return false
	}

	if err = sch.Validate(v); err != nil {
		log.Printf("%#v", err.Error())
		return false
	}

	return true
}

// For validation based on an external URL
func (e *Event) isValidByUrl(data []byte, schemaPath string) bool {
	reader, err := httploader.Load(schemaPath)
	if err != nil {
		log.Printf("%#v", err.Error())
		return false
	}

	compiler := jsonschema.NewCompiler()
	err = compiler.AddResource("schema.json", reader)
	if err != nil {
		log.Printf("%#v", err.Error())
		return false
	}
	sch, err := compiler.Compile("schema.json")
	if err != nil {
		log.Printf("%#v", err)
	}

	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		log.Printf("%#v", err.Error())
		return false
	}

	if err = sch.Validate(v); err != nil {
		log.Printf("%#v", err.Error())
		return false
	}

	return true
}
