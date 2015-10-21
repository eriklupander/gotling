package main

type StatFrame struct {
	Time int64 `json:"time"`
    Latency int64 `json:"reqs"`
}

type HttpReqResult struct {
	Latency int64
	Size int
	Status int
}
