package models

type ChangeLimitRequest struct {
	EndPoint  string `json:"endpoint"`
	Ratelimit int    `json:"ratelimit"`
}
