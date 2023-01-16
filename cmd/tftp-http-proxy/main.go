package main

import (
	"flag"
	"github.com/pin/tftp/v3"
	log "github.com/sirupsen/logrus"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"tftp-http-proxy/internal/handlers"
	"time"
)

func main() {
	var (
		host = flag.String("host", "0.0.0.0",
			"IP address of the host to listen on, 0.0.0.0 by default")
		port = flag.Int("port", 69,
			"port to listen on for incoming connections, 69 by default")
		baseUrl = flag.String("base-url", "",
			"The URL base to proxied request")
		logLevel = flag.String("log-level", "info",
			"The log level, use error,warning,info,debug. Info by defaul")
	)
	flag.Parse()

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.WithError(err).Warn("Failed to change log level to %s", *logLevel)
		level = log.InfoLevel
	}
	log.SetLevel(level)

	if *baseUrl == "" {
		log.Error("Missing URL base, specify the --baseUrl parameter")
		os.Exit(1)
	}

	_, err = url.Parse(*baseUrl)
	if err != nil {
		log.WithError(err).Error("Invalid URL base")
		os.Exit(1)
	}

	readHandler, err := handlers.NewHttpProxyGetHandler(*baseUrl, 5*time.Second)
	if err != nil {
		log.WithError(err).Error("Failed setting up proxy handler")
		os.Exit(1)
	}

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	tftpServer := tftp.NewServer(readHandler.Handler, nil)
	addr := net.JoinHostPort(*host, strconv.Itoa(*port))
	log.Infof("Starting TFTP server on %s", addr)

	go func() {
		err := tftpServer.ListenAndServe(addr)
		if err != nil {
			log.WithError(err).Error("Failed to start TFTP server")
			os.Exit(1)
		}
	}()

	<-sigChan
	log.Infof("Shutdown TFTP server on %s", addr)
	tftpServer.Shutdown()
}
