package event

import "time"

type BalanceUpdated struct {
	Name    string
	Payload interface{}
}

func NewBalanceUpdated() *BalanceUpdated {
	return &BalanceUpdated{
		Name: "BalanceUpdated",
	}
}

func (ba *BalanceUpdated) GetName() string {
	return ba.Name
}

func (ba *BalanceUpdated) GetPayload() interface{} {
	return ba.Payload
}

func (ba *BalanceUpdated) GetDateTime() time.Time {
	return time.Now()
}

func (ba *BalanceUpdated) SetPayload(payload interface{}) {
	ba.Payload = payload
}
