package pubsub

import (
	"encoding/json"
	"github.com/manjada/com/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"testing"
)

type Data struct {
	Text string
}

func TestPubSubRabbitMq_Send(t *testing.T) {
	type args struct {
		data PubSubData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test", args: struct{ data PubSubData }{data: PubSubData{
			Name:       "test",
			Exchange:   "test",
			RoutingKey: "",
			Mandatory:  false,
			Immediate:  false,
			Body:       Data{Text: "Test"},
		}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewRabbitMq()
			if err = r.Send(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPubSubRabbitMq_Receive(t *testing.T) {
	type args struct {
		exchange string
	}
	tests := []struct {
		name    string
		args    args
		want    <-chan amqp.Delivery
		wantErr bool
	}{
		{
			name:    "TEST_RECEIVE",
			args:    args{"test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewRabbitMq()
			got, err := r.Receive(tt.args.exchange)
			if (err != nil) != tt.wantErr {
				t.Errorf("Receive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for data := range got {
				m := Data{}
				err := json.Unmarshal(data.Body, &m)
				if err != nil {
					config.Error(err)
				}
				data.Ack(false)
			}
			/*go func() {
				for data := range got {
					m := Data{}
					err := json.Unmarshal(data.Body, &m)
					if err != nil {
						Error(err)
					}
					data.Ack(false)
				}
			}()*/

		})
	}
}
