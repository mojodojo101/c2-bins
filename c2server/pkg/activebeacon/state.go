package activebeacon

import (
	"context"
)

type State interface {
	Update(ctx context.Context, lineance int64) error
}
