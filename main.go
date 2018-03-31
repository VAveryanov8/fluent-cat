package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"
)

func main() {
	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)

	buffer := bufio.NewReader(os.Stdin)

	host := flag.String("h", "localhost", "fluentd host name")
	port := flag.Int("p", 24224, "fluentd port")
	asyncConnect := flag.Bool("async-connect", true, "reconnect if host unavailable")
	tag := flag.String("tag", "default", "tag")

	flag.Parse()

	client, err := fluent.New(fluent.Config{
		FluentHost:   *host,
		FluentPort:   *port,
		AsyncConnect: *asyncConnect,
	})
	if err != nil {
		logger.Fatalf("could not create fluentd client: %v\n", err)
	}

	for {
		line, _, err := buffer.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Printf("could not read line from stdin: %v\n", err)
			continue
		}
		err = client.EncodeAndPostData(*tag, time.Now(), line)
		if err != nil {
			logger.Printf("could not send event: %s\n", err)
		}
	}
}
