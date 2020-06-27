package job

import (
	"log"
	"time"
)

func Run(d time.Duration, name string, do func() error) *time.Ticker {
	log.Printf("(INFO) job %s started with repeat every: %v", name, d)
	ticker := time.NewTicker(d)
	go func() {
		for {
			if err := do(); err != nil {
				log.Printf("(ERR) job %s: error: %v", name, err)
			}
			if _, ok := <-ticker.C; !ok {
				log.Printf("(INFO) job %s: stopped", name)
				return
			}
		}
	}()
	return ticker
}
