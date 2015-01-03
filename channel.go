package pubnub

// Channel is used for established subscription connections
type Channel struct {
	Name      string
	MessageCh chan Message
	ErrorCh   chan error
	pr        *PubNubRequest
}

// NewChannel opens a new subscription channel
func (p *PubNubClient) NewChannel(name string) (*Channel, error) {
	if name == "" {
		return nil, ErrChannelNotSet
	}

	channel := &Channel{
		Name:      name,
		MessageCh: make(chan Message),
		ErrorCh:   make(chan error),
	}

	// channel not found
	pr := NewPubNubRequest(name, channel.MessageCh, channel.ErrorCh)
	channel.pr = pr

	go pr.handleResponse()

	// timetoken parameter is not sent for now
	go p.pub.Subscribe(name, "", pr.successCh, false, pr.errorCh)

	if err := pr.Do(); err != nil {
		return nil, err
	}

	return channel, nil
}

func (c *Channel) Close() {
	c.pr.Close()
	close(c.MessageCh)
	close(c.ErrorCh)
}
