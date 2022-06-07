package main

import (
	"go-torrent/src/Torrent"
	"log"
)

func main() {
	tf, err := Torrent.Open("D:\\GolandProject\\go-torrent\\src\\debian-iso.torrent")
	if err != nil {
		log.Fatal(err)
	}

	err = tf.DownloadToFile("D:\\GolandProject\\torrent-client")
	if err != nil {
		log.Fatal(err)
	}
}
