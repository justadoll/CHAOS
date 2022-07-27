package audio

import "context"

type Service interface {
	Record(ctx context.Context, address string, raw_seconds string) (string, error)
}
