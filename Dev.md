# Dev 

## Thanks to 

https://www.zupzup.org/go-http-file-upload-download/index.html
https://github.com/mayth/go-simple-upload-server/blob/v2/pkg/server.go
https://github.com/mjpclab/go-http-file-server
https://developer20.com/add-header-to-every-request-in-go/
https://lets-go.alexedwards.net/sample/02.04-customizing-http-headers.html


Problem breaking them down to what need to be done
- Temp to normal files https://gobyexample.com/temporary-files-and-directories - DONE

- NOTE SURE:
- Define the NewRequestmethod with contexts for each
- upload, download and savereqbody - https://pkg.go.dev/net/http#example-FileServer

- HEADERS
https://surajincloud.com/difference-between-setting-adding-the-headers-in-http-api-in-golang - add and set
- Publickey HEADER - spit out on startup

- Write into files and read file for /download 
- Would kind of like json because why not - https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

- Contexts: https://medium.com/@jamal.kaksouri/the-complete-guide-to-context-in-golang-efficient-concurrency-management-43d722f6eaea
- Follow the best practices 
- Timeout for connect - https://gobyexample.com/context
- Deadline for upload and requestbody
- Progating context with go routines
- canceling
- read - https://www.digitalocean.com/community/tutorials/how-to-use-contexts-in-go


- Make sure I am good https://www.digitalocean.com/community/tutorials/how-to-make-an-http-server-in-go

- Shutdown server https://pkg.go.dev/net/http#example-Server.Shutdown

- Final Writeup all my golang knownledge from this into Archive

Check Authentication Header copypasta code
```go
r.Header.Add("", "")

	if r.Header.Get("") != "" {
		shell = r.Header.Get("")
	} else {
		err := fmt.Errorf("...")
		checkError(err, 0)
	}

```




```go
// CLI -> if http else https -> Done - just simple done project - below is just a map of functions - see TODO idiot
// main -> handleArgs -> main
// switch on sortedArray -> subchecks on sortedArray size
// Either: http or https server

// TODO List TODO
// HTTP server
// Error ids and code - id for where in the source for no lost in the src and code for switch case fatal or not

// TLS is not that much more 
// TODO AFTER THE ABOVE IS COMPLETE built and works no excuses:
// Application end - start time printTotalRuntime()
// TLS - regex requred that make sense, validationTLS(), how validateTLS passes data to buildHTTPS()
// TODO AFTER THE ABOVE IS COMPLETE built and works no excuses:
// runHTTPSserver()

    log.Printf("About to listen on 8443. Go to https://127.0.0.1:8443/")
	err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
	log.Fatal(err)

// TODO AFTER THE ABOVE IS COMPLETE built and works no excuses:
// Add all the profession stuff
// - Make authentication actually work without any other packages
// - Hijacker? https://pkg.go.dev/net/http#example-Hijacker 
// - Colourful Banner!
// - DO NOT WORRY ABOUT nested regex -> string sorted args oneliners no (5||6)*2 additional variable declarations making that underreadable dense vertically and save some memory




```
