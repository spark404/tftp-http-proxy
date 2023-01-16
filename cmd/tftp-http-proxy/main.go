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
		listen = flag.String("listen", "0.0.0.0",
			"IP address to listen on for TFTP requests (default: \"0.0.0.0\")")
		port = flag.Int("port", 69,
			"The port to listen on for TFTP requests (default: 69)")
		baseUrl = flag.String("url", "",
			"URL to forward TFTP requests to")
		logLevel = flag.String("log-level", "info",
			"Sets the default log level (default: \"info\")")
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
	addr := net.JoinHostPort(*listen, strconv.Itoa(*port))
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
