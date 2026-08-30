package notificationgrpc

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/alexistamher/social-api-go/internal/domain/notification"
	notificationsv1 "github.com/alexistamher/social-api-go/internal/gen/proto/external/socialmedia/proto/socialmedia/notifications/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client notificationsv1.NotificationServiceClient
}

func NewClient(addr string, plaintext bool) (*Client, error) {
	var creds credentials.TransportCredentials
	if plaintext {
		creds = insecure.NewCredentials()
	} else {
		creds = credentials.NewTLS(&tls.Config{})
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("notificationgrpc: dial %q: %w", addr, err)
	}

	return &Client{
		conn:   conn,
		client: notificationsv1.NewNotificationServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) NotifyReaction(ctx context.Context, evt notification.ReactionEvent) error {
	content := &notificationsv1.NotifyContent{
		TargetId:   evt.ReactionID,
		TargetType: notificationsv1.TargetType_TARGET_TYPE_REACTION,
		ActionType: actionType(evt.Action),
	}
	if evt.ParentTargetID != "" {
		parent := evt.ParentTargetID
		content.ParentTargetId = &parent
	}

	_, err := c.client.Notify(ctx, &notificationsv1.NotifyRequest{
		ConnectionId: evt.ConnectionID,
		Content:      content,
	})
	if err != nil {
		return fmt.Errorf("notificationgrpc: notify: %w", err)
	}
	return nil
}

func actionType(a notification.Action) notificationsv1.ActionType {
	switch a {
	case notification.ActionAdd:
		return notificationsv1.ActionType_ACTION_TYPE_ADD
	case notification.ActionUpdate:
		return notificationsv1.ActionType_ACTION_TYPE_UPDATE
	case notification.ActionDelete:
		return notificationsv1.ActionType_ACTION_TYPE_DELETE
	default:
		return notificationsv1.ActionType_ACTION_TYPE_UNSPECIFIED
	}
}
