package api

type ClusterOperation struct {
	ClusterId     int64 `json:"clusterId"`
	OutputSubject string `json:"outputSubject"` // output nats subject where cluster outputs are sent
}

type ClusterCreateResponse struct {
	OutputChannel string
}

type TokenForm struct {
	Token          string `form:"jwt" binding:"Required"`
}
