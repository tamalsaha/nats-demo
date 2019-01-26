package api

type ClusterInfo struct {
	ClusterId int64
	OutputChannel string
}

type ClusterCreateResponse struct {
	OutputChannel string
}
