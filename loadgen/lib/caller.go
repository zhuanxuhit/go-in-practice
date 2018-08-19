package lib

import "time"

type Caller interface {
	BuildReq() RawReq
	Call(req []byte, timeoutNS time.Duration)([]byte, error)
	CheckResp(req RawReq, resp RawResp) *CallResult
}