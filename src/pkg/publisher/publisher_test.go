package publisher_test

import (
	"hms/gateway/pkg/publisher"
	"strconv"
	"sync"
	"testing"
)

type mockSubscriber struct {
	isClose    bool
	testNotify *func(string)
	name       string
}

func (s *mockSubscriber) Notify(msg interface{}) {
	(*s.testNotify)(msg.(string))
}
func (s *mockSubscriber) Disable() {
	s.isClose = true
}
func (s *mockSubscriber) Name() string {
	return s.name
}

func TestPublisher(t *testing.T) {
	pub := publisher.NewPublisher()

	t.Run("PublishMessage", func(t *testing.T) {
		wg := sync.WaitGroup{}

		cntSubs := 100
		wg.Add(cntSubs)

		msg := "test msg"
		cntSuccessDeliveries := 0
		testFunNotify := func(s string) {
			if msg != s {
				t.Errorf("expected:%s got:%s", msg, s)
			}
			cntSuccessDeliveries++
			wg.Done()
		}

		for i := 0; i < cntSubs; i++ {
			sub := mockSubscriber{
				isClose:    false,
				testNotify: &testFunNotify,
				name:       strconv.Itoa(i),
			}

			pub.AddSubscriber(&sub)
			defer func() {
				pub.RemoveSubscribe(&sub)
			}()
		}

		pub.PublishMessage(msg)

		wg.Wait()

		if cntSubs != cntSuccessDeliveries {
			t.Errorf("expected delivered messages:%d got:%d", cntSubs, cntSuccessDeliveries)
		}
	})
}
