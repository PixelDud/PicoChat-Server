/*
	PicoChat server by PixelDud

	Based on Hessery's original server made with GameMaker.
*/

package main

import (
	b32 "encoding/base32"
	"fmt"
	"net/http"
	"regexp"
	s "strings"
)

var chatLog []string // Chatlog array
var shortLog string  // String of messages seperated by a hyphen that gets sent when log is requested
var lastValidation string
var r, _ = regexp.Compile(`^/0.\d*/`)

// Function for encoding messages to Base32
func encodeMsg(msg string) string {
	out := b32.StdEncoding.EncodeToString([]byte(msg))
	return out
}

// Function for decoding messages from Base32
/* Unused at this time.
func decodeMsg(msg string) string {
	data, err := b32.StdEncoding.DecodeString(s.ToUpper(msg))
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(data)
}
*/

// Function for handling server requests
func handleRequest(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Last-Modified", "Wed, 12 Aug 1998 15:03:50 GMT")
	// Regex validator checks for decimal starting with zero of any length
	request := fmt.Sprintf("%s", req.URL)
	validator := r.FindString(request)
	// If validator isn't empty and doesn't match the last validation string, we move forward.
	if validator != "" && validator != lastValidation {
		request = s.Replace(request, validator, "", 1)
		lastValidation = validator
	} else {
		return
	}

	// Chech if the request's intention
	if s.HasPrefix(request, "%2bget") { // If it's trying to get the chatlog, send it.
		fmt.Fprint(w, shortLog)
		return
	} else {
		/*
			If it's doing anything else then we treat it as a message...
			which isn't the best way to go about this, but it works for the time being.
			I'm likely to add some other method of determining if it's a message to send or not.
		*/
		request = s.ToUpper(request) // Not necessary, but keeps the log uniform in terms of casing.
		// Chatlog is currently limited to 20 messages, client has no scrollback.
		if len(chatLog) < 20 {
			chatLog = append(chatLog, request)
			shortLog = s.Join([]string{shortLog, request + "-"}, "")
		} else { // When log is going to excede limit, we remove the first value and rewrite the shortlog.
			for i := 1; i < len(chatLog); i++ {
				chatLog[i-1] = chatLog[i]
			}
			chatLog[len(chatLog)-1] = request
			shortLog = chatLog[0] + "-"
			for i := 1; i < len(chatLog); i++ {
				shortLog = s.Join([]string{shortLog, chatLog[i] + "-"}, "")
			}
		}
		fmt.Fprint(w, shortLog)
		return
	}
}

// Initialize chat log with server start message and store in shortlog.
func init() {
	chatLog = make([]string, 1)
	chatLog[0] = encodeMsg("> Server started!\n")
	shortLog = chatLog[0] + "-"
}

// Where all the business happens.
func main() {
	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(":80", nil)
}
