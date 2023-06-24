package messaging

import (
	"bytes"
	"context"
	"encoding/json"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pkg/errors"
)

func New(
	ctx context.Context,
	selfHost host.Host,
	privateKey crypto.PrivKey,
	topic *pubsub.Topic,
	subscription *pubsub.Subscription,
) Service {
	return &service{
		ctx:          ctx,
		selfHost:     selfHost,
		privateKey:   privateKey,
		topic:        topic,
		subscription: subscription,
	}
}

var _ Service = (*service)(nil)

type service struct {
	ctx context.Context

	selfHost   host.Host
	privateKey crypto.PrivKey

	topic        *pubsub.Topic
	subscription *pubsub.Subscription
}

func (svc service) Send(to peer.ID, content []byte) error {
	publicKeyFrom, _ := svc.selfHost.ID().ExtractPublicKey()
	publicKeyTo, _ := to.ExtractPublicKey()

	body := payloadBody{
		FromID:      svc.selfHost.ID(),
		ToID:        to,
		FromContent: encode(content, publicKeyFrom),
		ToContent:   encode(content, publicKeyTo),
	}

	bodyData, _ := json.Marshal(body)

	if publishErr := svc.topic.Publish(svc.ctx, bodyData); publishErr != nil {
		return publishErr
	}

	return nil
}

func (svc service) GetNext(ctx context.Context) (Income, error) {
	var (
		message    *pubsub.Message
		messageErr error
		body       payloadBody
	)

	for {
		message, messageErr = svc.subscription.Next(ctx)
		if messageErr != nil {
			if messageErr == context.Canceled {
				return Income{}, messageErr
			}

			return Income{}, errors.Wrap(messageErr, "extracting message failed")
		}

		_ = json.Unmarshal(message.Data, &body)

		if body.ToID != svc.selfHost.ID() {
			continue
		}

		break
	}

	// decode data for self
	decodedData := decode(body.ToContent, svc.privateKey)

	// encode again with from-publicKey to validate signature
	pubKeyFrom, _ := message.ReceivedFrom.ExtractPublicKey()
	checkContent := encode(decodedData, pubKeyFrom)
	if bytes.Compare(checkContent, body.FromContent) != 0 {
		return Income{}, ErrUntrustedContentReceived
	}

	income := Income{
		Sent:    message.ReceivedFrom,
		Content: decodedData,
	}

	return income, nil
}

func (svc service) Close() error {
	svc.subscription.Cancel()
	return nil
}

func encode(b []byte, k crypto.PubKey) []byte {
	// TODO: implement me
	return b
}

func decode(b []byte, k crypto.PrivKey) []byte {
	// TODO: implement me
	return b
}
