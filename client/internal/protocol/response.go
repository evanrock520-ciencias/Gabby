package protocol

type Response struct {
	Operation TypeC2S `json:"operation,omitempty"`
	Result    Result  `json:"result,omitempty"`
	Extra     string  `json:"extra,omitempty"`
}
