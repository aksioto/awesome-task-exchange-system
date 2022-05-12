package rabbitmq

type Receiver func(body []byte)

//COMMON
type Event struct {
	EventID      string                 `json:"event_id"`
	EventVersion int                    `json:"event_version"`
	EventName    string                 `json:"event_name"`
	EventTime    int64                  `json:"event_time"`
	Producer     string                 `json:"producer"`
	Data         map[string]interface{} `json:"data"`
}
