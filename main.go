package main

import (
   // "flag"
    "time"
	"gopkg.in/cheggaaa/pb.v1"
)

type aaa struct {
	startTime int64
	bar       *pb.ProgressBar
    filename    string
    chuckList []string
}

func main() {
	handler := &aaa{
        startTime: time.Now().Unix(),
    }
	handler.askForFN()
	handler.parseM3U8()
	handler.start()
}
