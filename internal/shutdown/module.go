package shutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

// Module - init module.
type Module struct {
	Ctx context.Context

	wg    sync.WaitGroup
	count int32
}

// SafeRun starts a process, which must be completed before shutting down the system.
func (m *Module) SafeRun() {
	m.wg.Add(1)
	atomic.AddInt32(&m.count, 1)
}

// SafeComplete ends the process registered via SafeRun.
func (m *Module) SafeComplete() {
	m.wg.Done()
	atomic.AddInt32(&m.count, -1)
}

// Safe when receiving signals from the system (syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
// waits for all operations registered with SafeRun to complete SafeComplete.
// After terminating all handlers, terminate system execution.
func (m *Module) Safe() {
	go func() {
		ctx, stop := Subscribe(m.Ctx)
		<-ctx.Done()

		stop()

		log.Println("start safe stop application")

		count := atomic.LoadInt32(&m.count)

		if count > 0 {
			log.Println("waiting for active handlers to finish (active:% d)\n", count)
		}

		m.wg.Wait()
		os.Exit(0)
	}()
}

func Subscribe(ctx context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
}
