package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	addrs := []string{}
	timeout := 1 * time.Minute

	fs := flag.NewFlagSet("wait", flag.ContinueOnError)
	fs.Func("a", "Addresses to connect to [host:]port may be specified multilpe times", func(a string) error {
		if !strings.Contains(a, ":") {
			a = ":" + a
		}
		addrs = append(addrs, a)
		return nil
	})
	fs.DurationVar(&timeout, "t", timeout, "Timeout to wait")
	err := fs.Parse(os.Args[1:])
	if err != nil {
		os.Exit(2)
	}

	if len(addrs) == 0 {
		fs.Usage()
		os.Exit(2)
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	g, ctx := errgroup.WithContext(timeoutCtx)
	for _, addr := range addrs {
		addr := addr
		g.Go(func() error {
			for {
				log.Printf("[%s] waiting", addr)
				_, err = net.Dial("tcp", addr)
				if err == nil {
					break
				}
				select {
				case <-time.After(time.Second):
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			log.Printf("[%s] done", addr)
			return nil
		})
	}
	err = g.Wait()
	if err == context.DeadlineExceeded {
		log.Printf("timeout")
		os.Exit(3)
	} else if err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}
