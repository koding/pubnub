package pubnub

import "testing"

func TestGrantAccessWithGranter(t *testing.T) {
	cs := new(ClientSettings)
	cs.ID = "-1"
	cs.SubscribeKey = "subscribekey"
	cs.PublishKey = "publishkey"
	cs.SecretKey = "secretkey"

	ag := NewAccessGrant(NewAccessGrantOptions(), cs)

	a := new(AuthSettings)
	a.ChannelName = "testme"
	a.Token = "123"
	a.CanWrite = true
	a.CanRead = true
	err := ag.Grant(a)

	if err != nil {
		t.Errorf("Expected nil but got error while granting access: %s", err)
	}
}
