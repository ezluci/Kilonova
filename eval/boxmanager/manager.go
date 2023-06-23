package boxmanager

import (
	"context"

	"github.com/KiloProjects/kilonova"
	"github.com/KiloProjects/kilonova/eval"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

var _ eval.BoxScheduler = &BoxManager{}

// BoxManager manages a box with eval-based submissions
type BoxManager struct {
	dm kilonova.GraderStore

	numConcurrent int
	concSem       *semaphore.Weighted
	memSem        *semaphore.Weighted

	logger *zap.SugaredLogger

	availableIDs chan int
}

func (b *BoxManager) GetBox(ctx context.Context, memQuota int64) (eval.Sandbox, error) {
	if err := b.concSem.Acquire(ctx, 1); err != nil {
		return nil, err
	}
	if memQuota > 0 {
		if err := b.memSem.Acquire(ctx, memQuota); err != nil {
			return nil, err
		}
	}
	box, err := newBox(<-b.availableIDs, memQuota, b.logger)
	if err != nil {
		return nil, err
	}
	return box, nil
}

func (b *BoxManager) ReleaseBox(sb eval.Sandbox) {
	q := sb.MemoryQuota()
	if err := sb.Close(); err != nil {
		zap.S().Warnf("Could not release sandbox %d: %v", sb.GetID(), err)
	}
	b.availableIDs <- sb.GetID()
	b.memSem.Release(q)
	b.concSem.Release(1)
}

// Close waits for all boxes to finish running
func (b *BoxManager) Close(ctx context.Context) error {
	b.concSem.Acquire(ctx, int64(b.numConcurrent))
	close(b.availableIDs)
	return nil
}

// New creates a new box manager
func New(startingNumber int, count int, maxMemory int64, dm kilonova.GraderStore, logger *zap.SugaredLogger) (*BoxManager, error) {

	if startingNumber < 0 {
		startingNumber = 0
	}

	availableIDs := make(chan int, 3*count)
	for i := 1; i <= 2*count; i++ {
		availableIDs <- i + startingNumber
	}

	bm := &BoxManager{
		dm:            dm,
		concSem:       semaphore.NewWeighted(int64(count)),
		memSem:        semaphore.NewWeighted(maxMemory),
		availableIDs:  availableIDs,
		numConcurrent: count,

		logger: logger,
	}
	return bm, nil
}
