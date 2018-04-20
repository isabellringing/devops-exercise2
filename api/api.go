package api

type ClientRequest struct {
	SerialNumber uint32 `json:"serial_number"`
}

type ClientReply struct {
	Msg    string `json:"msg"`
	Secret string `json:"secret"`
}

type ClientError struct {
	Msg string `json:"msg"`
}
