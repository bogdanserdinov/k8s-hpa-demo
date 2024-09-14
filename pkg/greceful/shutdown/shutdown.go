package shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

// Graceful listens for either syscall.SIGINT or syscall.SIGTERM signals.
// Once one of these signals is received, it invokes the provided callback function
// to gracefully stop or clean up necessary components.
func Graceful(action func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		action()
	}()
}
