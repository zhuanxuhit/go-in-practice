package loadgen

import (
	"time"
	"go-in-practice/loadgen/lib"
	"context"
)

// 载荷发生器最基本的3个要素：
// 1. 超时时间
// 2. 载荷量 = 并发量 / 平均响应时间
// 3. 载荷持续时间
// 返回结果通过chan形式获取，chan是并发安全的
type myGenerator struct {
	caller      lib.Caller           // 调用器。
	timeoutNS   time.Duration        // 处理超时时间，单位：纳秒。
	lps         uint32               // 每秒载荷量。
	durationNS  time.Duration        // 负载持续时间，单位：纳秒。
	concurrency uint32               // 载荷并发量。
	tickets     lib.GoTickets        // Goroutine票池。
	ctx         context.Context      // 上下文。
	cancelFunc  context.CancelFunc   // 取消函数。
	callCount   int64                // 调用计数。
	status      uint32               // 状态。
	resultCh    chan *lib.CallResult // 调用结果通道。
}

func NewMyGenerator(
	caller lib.Caller,
	timeoutNS time.Duration,
	lps uint32,
	durationNS time.Duration) (lib.Generator, error) {

}
