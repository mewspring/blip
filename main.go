package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/speaker"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: blip [OPTIONS]... FILE")
	flag.PrintDefaults()
}

func main() {
	var dur time.Duration
	flag.DurationVar(&dur, "dur", time.Second/10, "duration between rebuffering of samples")
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	fr, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatalf("unable to open file; %v", err)
	}
	s, format, err := flac.Decode(fr)
	speaker.Init(format.SampleRate, format.SampleRate.N(dur))
	done := make(chan struct{})
	f := func() {
		close(done)
	}
	notify := beep.Callback(f)
	speaker.Play(beep.Seq(s, notify))
	select {
	case <-done:
		break
	}
}
