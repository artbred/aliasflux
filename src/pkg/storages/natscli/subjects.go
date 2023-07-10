package natscli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

type Subject string

func (s Subject) String() string {
	return string(s)
}

func PerformRequest(ctx context.Context, subject Subject, req interface{}, res interface{}) error {
	nc := Connection()

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error while marshaling BillingCreatePaymentLinkRequest: %v", err)
	}

	msg, err := nc.RequestWithContext(ctx, subject.String(), data)
	if err != nil {
		return fmt.Errorf("error while sending request to nats: %v", err)
	}

	logrus.Print(string(msg.Data))

	return json.Unmarshal(msg.Data, res)
}
