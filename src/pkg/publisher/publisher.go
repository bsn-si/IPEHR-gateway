package publisher

type Publisher struct {
	subscribers map[string]Subscriber
	addSubCh    chan Subscriber
	removeSubCh chan Subscriber
	msg         chan interface{}
	stop        chan struct{}
}

func (p *Publisher) AddSubscriber(sub Subscriber) {
	p.addSubCh <- sub
}

func (p *Publisher) RemoveSubscribe(sub Subscriber) {
	p.removeSubCh <- sub
}

func (p *Publisher) PublishMessage(msg interface{}) {
	p.msg <- msg
}

func (p *Publisher) Stop() {
	close(p.stop)
}

func (p *Publisher) start() {
	for {
		select {
		case sub := <-p.addSubCh:
			{
				p.subscribers[sub.Name()] = sub
			}
		case sub := <-p.removeSubCh:
			{
				if s, ok := p.subscribers[sub.Name()]; ok {
					s.Disable()
					delete(p.subscribers, sub.Name())
				}
			}
		case msg := <-p.msg:
			{
				for _, sub := range p.subscribers {
					sub.Notify(msg)
				}
			}
		case <-p.stop:
			{
				for _, sub := range p.subscribers {
					sub.Disable()
				}

				close(p.addSubCh)
				close(p.removeSubCh)
				close(p.msg)

				return
			}
		}
	}
}

func NewPublisher() *Publisher {
	em := Publisher{
		subscribers: map[string]Subscriber{},
		addSubCh:    make(chan Subscriber),
		removeSubCh: make(chan Subscriber),
		msg:         make(chan interface{}),
		stop:        make(chan struct{}),
	}
	go em.start()
	return &em
}
