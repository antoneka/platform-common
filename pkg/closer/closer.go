package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

var globalCloser = New()

// Add adds closure functions to closer.
func Add(f ...func() error) {
	globalCloser.Add(f...)
}

// Wait blocks until all closure functions are done.
func Wait() {
	globalCloser.Wait()
}

// CloseAll calls all closure functions.
func CloseAll() {
	globalCloser.CloseAll()
}

// Closer represents a closer that manages a collection of closure functions.
type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []func() error
}

// New returns new Closer.
// If []os.Signal is specified Closer will automatically call CloseAll
// when one of signals is received from OS.
func New(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()
	}
	return c
}

// Add adds closure functions to closer.
func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

// Wait blocks until all closure functions are done.
func (c *Closer) Wait() {
	<-c.done
}

// CloseAll calls all closure functions.
func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		// Call all closure funcs async.
		errs := make(chan error, len(funcs))
		for _, f := range funcs {
			go func(f func() error) {
				errs <- f()
			}(f)
		}

		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				log.Println("error returned from Closer")
			}
		}
	})
}
