package protocol

type Response struct {
	Operation TypeC2S
	Result    Result
	Extra     string
}
