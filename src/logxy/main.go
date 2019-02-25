package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)


type Logxy struct {
	forwardTo  string
	port       int
}


func NewLogxy(forwardTo string, port int) *Logxy {

	logxy := &Logxy{
		forwardTo:  forwardTo,
		port:       port,
	}

	return logxy
}


func (logxy *Logxy) Close() error {

	logxy.forwardTo = ""
	logxy.port = -1
	return nil
}


func (logxy *Logxy) Run() {

	addr := fmt.Sprintf(":%v",logxy.port)

	// setup endpoint handlers
	http.HandleFunc("/", logxy.requestHdl)

	// start up server
	log.Fatal(http.ListenAndServe(addr, nil))
}


// log requests to the reverse proxy's target
func logRequest(req *http.Request) (err error) {

	// read the request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return err
	}

	// close the body
	err = req.Body.Close()
	if err != nil {
		return err
	}

	// https://github.com/bechurch/reverse-proxy-demo/blob/master/main.go#L116
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// log the request
	log.Printf("request: %v %v %v", req.Method, req.URL.String(), string(body))

	// return no error
	return nil
}

// log responses from the reverse proxy's target
func logResponse(res *http.Response) (err error) {

	// read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// close the body
	err = res.Body.Close()
	if err != nil {
		return err
	}

	// https://github.com/bechurch/reverse-proxy-demo/blob/master/main.go#L116
	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// log the response
	log.Printf("response: %v", string(body))

	// return no error
	return nil
}


// Handle requests to the / endpoint
func (logxy *Logxy) requestHdl(w http.ResponseWriter, req *http.Request) {

	logRequest(req)

	// get the target url
	url, _ := url.Parse(logxy.forwardTo)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// setup handler to log responses from targer url
	proxy.ModifyResponse = logResponse

	// update the headers for ssl redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forward-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// run the reverse proxy
	proxy.ServeHTTP(w, req)
}


// helper function from https://bit.ly/2j4nNSs
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}


// helper function from https://bit.ly/2j4nNSs
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}


type cmdOptions struct {
	forwardTo  *string
	log        *string
	port       *int
}


func (opts *cmdOptions) String() string {
	return fmt.Sprintf("cmdOptions{ forwardTo: %v, log: %v, port: %v }", *opts.forwardTo, *opts.log, *opts.port)
}


func parseOptions() *cmdOptions {

	opts := &cmdOptions{}

	// command line arguments

	opts.forwardTo = flag.String(
		"forward-to",
		"http://selenium-hub:4444",
		"URL of the server to forward requests to")

	opts.log = flag.String(
		"log",
		"./logxy.log",
		"Path to write the log messages to")

	opts.port = flag.Int(
		"port",
		8744,
		"Port this reverse proxy server should listen on for HTTP requests")

	flag.Parse()

	return opts
}


func setupLogging(logFilename string) *os.File {

	// open a logfile
	logFile, err := os.OpenFile(logFilename,
		os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	// send logs to a logfile
	log.SetOutput(logFile)

	return logFile
}


func setupSignalCaptures(closers [](func() error)) {

	// open a channel that we can send signals to
	sigs := make(chan os.Signal, 1)

	// when we get a signal, notify the sigs channel
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// use a goroutine/thread to block until we get a signal
	// when a signal arrives on the sigs channel, call the cleanup methods
	go func() {

		// block, waiting for a signal
		<-sigs

		// cleanup, close file handles, ...
		cleanup(closers)

		//exit
		os.Exit(0)
	}()
}


func cleanup(closers [](func() error)) {

	log.Printf("Cleaning up before exiting...")

	// run the closers
	for _, closer := range closers {
		closer()
	}
}


// create a Logxy that accepts requests, logs them,
// and forwards them on to the target server

func main() {

	opts := parseOptions()

	// setup logging to stdout and a logfile
	logFile := setupLogging(*opts.log)
	defer logFile.Close()

	log.Printf("Starting Logxy")
	log.Printf("args = %v", os.Args)
	log.Printf("opts = %v", opts)

	// create a new Logxy
	logxy := NewLogxy(
		*opts.forwardTo,
		*opts.port,
	)
	defer logxy.Close()

	// register things to close when cleaning up
	// close Logxy server
	// close the logfile
	closers := [](func() error){logxy.Close, logFile.Close}

	// capture signals
	setupSignalCaptures(closers)

	// start the reverse proxy server
	logxy.Run()
}

