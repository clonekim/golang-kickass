package main

import (
	"flag"
	"log"
	"os"
	"time"
)

const IDLE_TIME = 100 * time.Millisecond

func parse() *os.File {

	filename := flag.String("filename", "", "select filename to watch and run")
	flag.Parse()

	if *filename == "" {
		panic("filename is required")
	}

	file, err := os.Open(*filename)

	if err != nil {
		log.Println(err)

		file, err := os.Create(*filename)

		if err != nil {
			panic(err)
		}

		return file

	} else {
		return file
	}

}

func main() {

	file := parse()
	defer file.Close()

	last := time.Now()

	for {
		changed, updateT := watchDetect(file, last)

		if changed {
			log.Println("file modified")
			last = updateT
		}

		time.Sleep(IDLE_TIME)
	}

}

func watchDetect(file *os.File, t time.Time) (bool, time.Time) {

	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}

	//t보다 최신 날짜인가?
	return fileInfo.ModTime().After(t), fileInfo.ModTime()
}
