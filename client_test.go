package pubnub

import (
	"sync"
	"testing"
	"time"
)

func createClient(id string) *PubNubClient {
	cs := newClientSettings(id)

	return NewPubNubClient(cs)
}

func createMessage() map[string]string {
	return map[string]string{
		"eventName": "messageAdded",
		"body":      "testing all together",
	}
}

func TestSubscribe(t *testing.T) {
	pc := createClient("tester")
	defer pc.Close()
	_, err := pc.Subscribe("")

	if err != ErrChannelNotSet {
		t.Errorf("Expected channel is not set error but got %s", err)
	}

	channels := "testme"
	_, err = pc.Subscribe(channels)
	if err != nil {
		t.Errorf("Expected nil but got error while subscribing: %s", err)
	}
}

func TestPublish(t *testing.T) {
	pc := createClient("tester")
	defer pc.Close()

	message := createMessage()
	err := pc.Push("tester1", message)
	if err != nil {
		t.Errorf("Expected nil but got error while publishing: %s", err)
	}
}

func TestMessageReception(t *testing.T) {
	sender := createClient("tester")
	receiver := createClient("receiver")

	defer func() {
		sender.Close()
		receiver.Close()
	}()

	// these subscriptions are added for checking if messages are received correctly
	// when we are subscribed to more than one channels
	receiver.Subscribe("sueme")

	channel, err := receiver.Subscribe("testme")
	if err != nil {
		t.Errorf("Expected nil but got error while subscribing: %s", err)
		t.FailNow()
	}

	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		select {
		case msg := <-channel.MessageCh:
			body, ok := msg.Body.(map[string]interface{})
			if !ok {
				t.Errorf("Wrong message body type")
				t.FailNow()
			}

			val, ok := body["eventName"]
			if !ok {
				t.Error("'eventName' field is expected in message but not found")
				t.Fail()
			} else {
				if val != "messageAdded" {
					t.Errorf("Expected messageAdded event in message but got %s", val)
				}
			}

			val, ok = body["body"]
			if !ok {
				t.Error("'body' field is expected in message but not found")
				t.Fail()
			} else {
				if val != "testing all together" {
					t.Errorf("Expected 'testing all together' as message body but got %s", val)
				}
			}

		case <-time.After(5 * time.Second):
			t.Errorf("Expected message but it is timedout")
			t.FailNow()
		}
	}()

	err = sender.Push("testme", createMessage())
	if err != nil {
		t.Errorf("Expected nil but got error while publishing: %s", err)
	}
	// TODO when we concurrently push message and subscribe to a channel, pushed message
	// is sometimes lost, due to reconnection
	// receiver.Subscribe("padme")

	wg.Wait()

}

func TestGrantAccess(t *testing.T) {
	gc := createClient("tester")
	defer func() {
		gc.Close()
	}()

	a := new(AuthSettings)
	a.ChannelName = "testme"
	a.Token = "123"
	a.CanWrite = true
	a.CanRead = true
	err := gc.Grant(a)

	if err != nil {
		t.Errorf("Expected nil but got error while granting access: %s", err)
	}

}
