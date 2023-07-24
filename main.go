package main

import (
	"flag"
	"fmt"
	"net/http"
)

type COLOR string

const (
	BLACK   COLOR = `30`
	RED     COLOR = `31`
	GREEN   COLOR = `32`
	YELLOW  COLOR = `33`
	BLUE    COLOR = `34`
	MAGENTA COLOR = `35`
	CYAN    COLOR = `36`
	WHITE   COLOR = `37`
)

func colorPrinter(color COLOR, format string, values ...any) {
	format = "\x1b[" + string(color) + "m" + format + "\x1b[0m"
	fmt.Printf(format, values...)
}

func main() {
	url := flag.String("u", "", "url to parse")
	flag.Parse()
	colorPrinter(BLUE, "Checking... URL : %s\n", *url)
	fmt.Println()
	resp, err := http.Get(*url)
	if err != nil {
		panic(err)
	}
	
	colorPrinter(YELLOW, "[Headers]\n")
	for key, values := range resp.Header {
		for _, value := range values {
			colorPrinter(MAGENTA, "	%-35s : ", key)
			colorPrinter(WHITE, "%s\n", value)
		}
	}
	fmt.Println()

	// Check Secure Cookies
	secureCookie := false
	colorPrinter(YELLOW, "[Cookie]\n")
	for _, cookie := range resp.Cookies() {
		if cookie.Secure {
			fmt.Printf("	%s\n", cookie.String())
			secureCookie = true
		}
	}
	fmt.Println()

	hsts := resp.Header.Get("Strict-Transport-Security")
	if hsts != "" {
		colorPrinter(GREEN, "[OK]:HSTS is set: %s\n", hsts)
	} else {
		colorPrinter(RED, "[CRITICAL]:HSTS is not set.\n")
	}
	if secureCookie {
		colorPrinter(GREEN, "[OK]:Secure attribute is set on one or more cookies.\n")
	} else {
		colorPrinter(RED, "[CRITICAL]:Secure attribute is not set on any of the cookies.\n")
	}
}
