package socket

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/subscriptions"
	"golang.org/x/sync/errgroup"
)

// SendNotify 向pb前端发送通知
func SendNotify(app core.App, subscription string, data any) error {
	rawData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	message := subscriptions.Message{
		Name: subscription,
		Data: rawData,
	}

	group := new(errgroup.Group)

	chunks := app.SubscriptionsBroker().Clients()

	for _, client := range chunks {

		if !client.HasSubscription(subscription) {
			continue
		}

		client.Send(message)
	}

	return group.Wait()
}
