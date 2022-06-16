package main

import (
	"fmt"
	"time"
	"os"
	"os/signal"
	"syscall"
)

type sigrecv_chan chan os.Signal
var sgchn sigrecv_chan
func (s *sigrecv_chan) sighandle() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,
		syscall.SIGINT,
		syscall.SIGTERM)

	done := make(chan os.Signal, 1)

	go func() {
		sig := <-sigs
		done <- sig
	}()
	*s = done
}

func (s sigrecv_chan) checksignal() (bool) {
	select {
	case sig := <-s:
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, sig)
		return true
	default:
	}
	return false
}

func sleep(second int64) {
	fmt.Fprintf(os.Stderr, "%s Sleep: %d", time.Now().Format("15:04:05"), second)
	start := time.Now()
	startunix := start.Unix()
	lastunix := startunix + int64(second)
	
	for second > 0 {
		slp := second
		if slp > sleepdot {
			slp = sleepdot
		}
		
		//time.Sleep(time.Duration(slp) * time.Second)
		for i := int64(0); i <= slp; i++ {
			if sgchn.checksignal() {
				print_id()
				os.Exit(3)
			}
			if i < slp {
				time.Sleep(time.Second)
			}
		}
		fmt.Fprintf(os.Stderr, ".")

		now := time.Now()
		nowunix := now.Unix()
		second = lastunix - nowunix
		if second < -10 {
			fmt.Fprintf(os.Stderr, "oversleep %s\n", now.Format("15:04:05"))
			// print_id()
			// os.Exit(0)
		}
		if second <= 0 {
			break
		}
	}
}
