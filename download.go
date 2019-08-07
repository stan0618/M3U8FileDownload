package main

import (
	"bufio"
	"github.com/grafov/m3u8"
    "errors"
	"strconv"
    "fmt"
	"log"
    "os"
	"time"
	"net/http"
	"io"
	"github.com/korovkin/limiter"
	"gopkg.in/cheggaaa/pb.v1"
)

func (h *aaa) askForFN() {
    var filename string

    filename = h.filename
    for filename == "" {
        fmt.Println("Please type anime sn id: ")
        fmt.Scanln(&filename)
    }
	
	// Check chunk exist or not
    fi, err := os.Stat(filename)
    if err != nil && fi.Size() == 0 {
        isErr("Extract filename failed -", errors.New("filename not found"))
    }
    h.filename = filename
}

func (h *aaa) parseM3U8() {
    f, err := os.Open(h.filename)
    isErr("Failed to read m3u8 playlist -", err)

    pl, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
    isErr("Parse m3u8 playlist failed -", err)

    switch listType {
    case m3u8.MEDIA:
        mediapl := pl.(*m3u8.MediaPlaylist)

        for _, chuck := range mediapl.Segments {
            if chuck != nil {
				fmt.Println(chuck.URI)
                h.chuckList = append(h.chuckList, chuck.URI)
            }
        }

        fmt.Println("All segments parsed! Download is starting...")
    }
}

func (h *aaa) start() {
    h.bar = pb.StartNew(len(h.chuckList))
    limit := limiter.NewConcurrencyLimiter(64)
	
	count := 0
    for _, url := range h.chuckList {
		count++
        part := url
        limit.Execute(func() {
            for {
                if h.downloadChunk(part, count) {
                    h.bar.Increment()
                    break
                }
            }
        })
    }
    limit.Wait()

    h.bar.Finish()
}

func (h *aaa) downloadChunk(chuckUrl string, count int) bool {
	t := strconv.Itoa(count)
    filename := "tmpMV"+t+".ts"

    // Check chunk exist or not
    fi, err := os.Stat(filename)
    if err == nil && fi.Size() != 0 {
        return true
    }

    // Create a chunk file
    out, err := os.Create(filename)
    isErr("Create "+filename+" failed -", err)

    req, err := http.NewRequest("GET", chuckUrl, nil)
    isErr("Create request failed - ", err)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        fmt.Println("Download "+filename+" file failed -", err)
        fmt.Println("Retrying -", filename)
        out.Close()
        os.Remove(filename)
        time.Sleep(500 * time.Millisecond)
        return false
    }

    _, err = io.Copy(out, resp.Body)
    if err != nil {
        fmt.Println(filename+" save failed -", err)
        fmt.Println("Retrying -", filename)
        out.Close()
        os.Remove(filename)
        time.Sleep(500 * time.Millisecond)
        return false
    }
    out.Close()
    return true
}

func isErr(msg string, err error) {
    if err != nil {
        f, e := os.OpenFile("error.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
        if e != nil {
            log.Fatal(err.Error())
        }

        log.SetOutput(f)
        msg := msg + " " + err.Error()
        fmt.Println(msg)
        log.Fatal(msg)
    }
}

