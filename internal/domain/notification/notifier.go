package notification

import "context"

type Action int

const (
	ActionAdd    Action = iota + 1 // entity was created
	ActionUpdate                   // entity was modified
	ActionDelete                   // entity was removed
)

func (a Action) String() string {
	switch a {
	case ActionAdd:
		return "add"
	case ActionUpdate:
		return "update"
	case ActionDelete:
		return "delete"
	default:
		return "unspecified"
	}
}

type ReactionEvent struct {
	ConnectionID   string
	ReactionID     string
	ParentTargetID string
	Action         Action
}

type Notifier interface {
	NotifyReaction(ctx context.Context, evt ReactionEvent) error
}

type NoopNotifier struct{}

func (NoopNotifier) NotifyReaction(context.Context, ReactionEvent) error { return nil }
