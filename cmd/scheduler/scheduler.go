package main

import (
	"fmt"
	"time"
)

func newScheduler(cfg *appConfig) (*scheduler, error) {
	sdl := scheduler{
		cfg: cfg,
	}
	return &sdl, nil
}

func (sdl *scheduler) Start() {

	for {
		fmt.Println("Num workers: ", len(sdl.cfg.workers))
		for _, w := range sdl.cfg.workers {
			fmt.Println("   ", w.Host, w.Port)
		}

		time.Sleep(time.Second * 10)
	}
}
