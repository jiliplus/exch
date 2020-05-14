package backtest

// NextIDFunc 返回的函数，可以生成连续的 ID
func NextIDFunc() func() int64 {
	id := int64(0)
	return func() int64 {
		id++
		return id
	}
}
