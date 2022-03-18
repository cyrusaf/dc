package dc

import (
	"context"
	"time"
)

type Updater interface {
	Update(ctx context.Context) error
}

type Poller struct {
	config   Updater
	interval time.Duration
	onError  func(err error)
	shutdown chan struct{}
}

func NewPoller(config Updater, interval time.Duration, onError func(err error)) Poller {
	return Poller{
		config:   config,
		interval: interval,
		shutdown: make(chan struct{}, 1),
	}
}

func (p Poller) Poll(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	for {
		select {
		case <-ticker.C:
			err := p.config.Update(ctx)
			if err != nil && p.onError != nil {
				p.onError(err)
			}
		case <-ctx.Done():
			return
		case <-p.shutdown:
			return
		}
	}
}

func (p Poller) Shutdown() {
	select {
	case p.shutdown <- struct{}{}:
	default:
	}
}
