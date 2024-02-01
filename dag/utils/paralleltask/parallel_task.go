package paralleltask

import (
	"context"
	"fmt"
	logging "github.com/ipfs/go-log/v2"
	"sync"
)

var log = logging.Logger("parallel-task")

type ParallelTask struct {
	entrySuccessQuorum int
	entryFailureQuorum int
	cancelOther        bool
	resCh              chan error
	parentCtx          context.Context
	cancelCtx          context.Context
	cancel             context.CancelFunc
	wg                 sync.WaitGroup
}

func NewParallelTask(ctx context.Context, entrySuccessQuorum, entryFailureQuorum int, cancelOtherGoroutine bool) *ParallelTask {
	cctx, cancel := context.WithCancel(ctx)
	return &ParallelTask{
		entrySuccessQuorum: entrySuccessQuorum,
		entryFailureQuorum: entryFailureQuorum,
		cancelOther:        cancelOtherGoroutine,
		resCh:              make(chan error),
		parentCtx:          ctx,
		cancelCtx:          cctx,
		cancel:             cancel,
	}
}

func (et *ParallelTask) Goroutine(f func(ctx context.Context) error) {
	et.wg.Add(1)
	go func() {
		defer func() {
			et.wg.Done()
			if rerr := recover(); rerr != nil {
				rerror := fmt.Errorf("catch panic :%v", rerr)
				log.Error(rerror)
			}
		}()
		ctx := et.parentCtx
		if et.cancelOther {
			ctx = et.cancelCtx
		}
		err := f(ctx)
		select {
		case <-et.cancelCtx.Done():
			return
		case et.resCh <- err:
		}
	}()
}

func (et *ParallelTask) Wait() error {
	defer et.clean()
	successCount := 0
	failureCount := 0
	for {
		select {
		case err := <-et.resCh:
			if err == nil {
				successCount += 1
				if successCount >= et.entrySuccessQuorum {
					return nil
				}
			} else {
				failureCount += 1
				if failureCount >= et.entryFailureQuorum {
					// return last error
					return err
				}
			}

		case <-et.cancelCtx.Done():
			return et.cancelCtx.Err()
		}
	}

}

func (et *ParallelTask) clean() {
	et.cancel()
	go func() {
		et.wg.Wait()
		close(et.resCh)
	}()
}
