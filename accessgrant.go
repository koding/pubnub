package pubnub

import (
	"sync"

	"github.com/pubnub/go/messaging"
)

type AccessGrant struct {
	pool sync.Pool
}

type AccessGrantOptions struct {
	ResumeOnReconnect bool
	SubscribeTimeout  int
	Origin            string
}

func NewAccessGrantOptions() *AccessGrantOptions {
	return &AccessGrantOptions{
		ResumeOnReconnect: true,
		SubscribeTimeout:  20,
		Origin:            "pubsub.pubnub.com",
	}
}

func NewAccessGrant(ao *AccessGrantOptions, cs *ClientSettings) *AccessGrant {
	messaging.SetResumeOnReconnect(ao.ResumeOnReconnect)
	messaging.SetSubscribeTimeout(uint16(ao.SubscribeTimeout))
	messaging.SetOrigin(ao.Origin)

	p := sync.Pool{
		New: func() interface{} {
			return NewPubNubClient(cs)
		},
	}

	return &AccessGrant{pool: p}
}

func (ag *AccessGrant) Grant(as *AuthSettings) error {
	client := ag.pool.Get().(*PubNubClient)
	defer ag.pool.Put(client)

	client.pub.SetAuthenticationKey(as.Token)

	return client.Grant(as)
}
