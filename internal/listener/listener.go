package listener

type ConfigQueue string

type Config struct {
	Queue *struct {
		Name   string            `json:"name"`
		Params map[string]string `json:"params,omitempty"`
	} `json:"queue,omitempty"`
	Logger *string `json:"logger,omitempty"`
}
