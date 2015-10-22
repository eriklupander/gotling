package main

type StatFrame struct {
	Time int64 `json:"time"`
    Latency int `json:"latency"`
	Reqs int `json:"reqs"`
}

type HttpReqResult struct {
	Latency int64
	Size int
	Status int
}
