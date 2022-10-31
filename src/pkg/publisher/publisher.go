package publisher

import "log"

type Publisher struct {
	subscribers []Subscriber
	addSubCh    chan Subscriber
	removeSubCh chan Subscriber
	inMsg       chan interface{}
	stop        chan struct{}

	addSubHandler    func(Subscriber)
	removeSubHandler func(Subscriber)
}

func (p *Publisher) AddSubscriber() chan<- Subscriber {
	return p.addSubCh
}
func (p *Publisher) RemoveSubscribe() chan<- Subscriber {
	return p.removeSubCh
}
func (p *Publisher) PublishMessage() chan<- interface{} {
	return p.inMsg
}
func (p *Publisher) Stop() {
	close(p.stop)
}

func (p *Publisher) onAddSubscriber(sub Subscriber) {
	if p.addSubHandler != nil {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic onAddSubscriber:%v", r)
			}
		}()

		p.addSubHandler(sub)
	}
}
func (p *Publisher) onRemoveSubscriber(sub Subscriber) {
	if p.removeSubHandler != nil {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic onRemoveSubscriber:%v", r)
			}
		}()

		p.removeSubHandler(sub)
	}
}

func (p *Publisher) start() {
	for {
		select {
		case sub := <-p.addSubCh:
			{
				p.subscribers = append(p.subscribers, sub)
				p.onAddSubscriber(sub)
			}
		case sub := <-p.removeSubCh:
			{
				for i, s := range p.subscribers {
					if sub == s {
						p.subscribers = append(p.subscribers[:i], p.subscribers[i+1:]...)
						s.Disable()
						p.onRemoveSubscriber(sub)
						break
					}
				}
			}
		case msg := <-p.inMsg:
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
				close(p.inMsg)

				return
			}
		}
	}
}

func NewPublisher() *Publisher {
	em := Publisher{
		addSubCh:    make(chan Subscriber),
		removeSubCh: make(chan Subscriber),
		inMsg:       make(chan interface{}),
		stop:        make(chan struct{}),
	}
	go em.start()
	return &em
}
