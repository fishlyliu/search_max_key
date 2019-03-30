# search_max_key

# 题目
有一个函数，行为如下 (Go 代码):
```
// need to call rand.Read(maxKey) during init

var maxKey = make([]byte, 256)

func init() {
	rand.Read(maxKey)
}

func Search(key []byte) []byte {
	time.Sleep(time.Millisecond*10)
	if bytes.Compare(key, maxKey) > 0 {
		return nil
	}
	return key
}
```

其中 maxKey 是未知的，需要通过多次调用 Search 函数得到最终的 maxKey 的值。每次调用 Search 函数都会有 10ms 的延迟。要求在多核多线程的条件下，充分利用所有计算资源，在最短的时间内计算得到 maxKey 的值。

# 分析
rand.Read随机生成256长度的byte数组maxKey。需求是从256长度的byte数组空间范围中找到该值。
单线程的话，应该用2分的方式来解决。若在多核多线程的条件下，容易考虑使用与核数相同的线程数量来处理。但要注意，题目限制比较时只能使用Search函数，而该函数有固定的time.Sleep 10ms行为，而time.Sleep本身不消耗cpu时间。

# 处理逻辑
1. 由于Search函数必须与完整的maxKey做字符串比较，所以整体思路是一个字符一个字符来决定；
2. 单次确定逻辑，可以开256个协程来同时度过10ms等待时间，故理论上总时间应该是10ms * 256；

# 测试结果，
在四核16g的机器上，cpu压力10%，耗时2.9秒，和预期比较接近。
可直接使用 go run step_search_td.go验证。

# 进一步思考
每次一个byte的比较，时间是256乘10ms。那如果每次两个byte，那么时间就应该是128乘10ms。
进而如果能256个字符一起比较，那么极限时间，应该就是10ms。
但是这里要考虑协程创建和调度的成本，每次两个byte为例，一次需要创建256*256个协程，调度上有巨大压力，成为性能瓶颈。在同样的四核16g机器上测试，cpu压力为230%，耗时为15s。

# 总结
该题本次解答，找到的最快处理方式耗时2.9秒。瓶颈在于协程调度消耗。
