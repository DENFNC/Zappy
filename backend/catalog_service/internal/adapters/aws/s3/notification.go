package s3client

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var NotifyConfigs []*types.NotificationConfiguration

type Notifyer struct {
	client *Client
}

type TypeEvent string

const (
	MethodPut TypeEvent = "PUT"
	// MethodGet TypeEvent = "GET"
)

func NewNotifyer(client *Client) *Notifyer {
	return &Notifyer{
		client: client,
	}
}

func (ntf *Notifyer) RegisterNewNotify(
	ctx context.Context,
	id, queueArn, bucket string,
	method TypeEvent,
) error {
	var qc types.QueueConfiguration

	switch method {
	case "PUT":
		qc = types.QueueConfiguration{
			Id:       aws.String(id),
			QueueArn: aws.String(queueArn),
			Events: []types.Event{
				types.EventS3ObjectCreatedPut,
			},
		}
	default:
		return errors.New("unknown type of event")
	}

	newNotif := &types.NotificationConfiguration{
		QueueConfigurations: []types.QueueConfiguration{qc},
	}

	_, err := ntf.client.API.PutBucketNotificationConfiguration(ctx,
		&s3.PutBucketNotificationConfigurationInput{
			Bucket:                    aws.String(bucket),
			NotificationConfiguration: newNotif,
		},
	)
	return err
}
