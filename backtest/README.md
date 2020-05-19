# backtest

根据 <https://github.com/ThreeDotsLabs/watermill/blob/482e08e8bb2a28b1c80dd4c00792fa8d23391f5e/pubsub/gochannel/pubsub.go#L194> 中的源代码

```go
func (g *GoChannel) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	// ....
	go func(s *subscriber, g *GoChannel) {
		select {
		case <-ctx.Done():
			// unblock
		case <-g.closing:
			// unblock
		}
		s.Close()
	    // ....
	}(s, g)
	// ....
	return s.outputChannel, nil
}
```

GoChannel.Close() 后，所有的 subscriber 也会被 close。
