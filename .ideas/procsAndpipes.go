//
cmd := exec.Command("command", "arg1", "arg2")

stdin, err := cmd.StdinPipe()
if err != nil {
	log.Fatal(err)
}
stdout, err := cmd.StdoutPipe()
if err != nil {
	log.Fatal(err)
}
stderr, err := cmd.StderrPipe()
if err != nil {
	log.Fatal(err)
}
err = cmd.Start()
if err != nil {
	log.Fatal(err)
}

err = cmd.Stop()
if err != nil {
	log.Fatal(err)
}

_, err = stdin.Write([]byte(""))
if err != nil {
	log.Fatal(err)
}

// Double check os.Pipe() for use channels over IPC for go routines in destination processes
r, w, err := os.Pipe()
if err != nil {
	log.Fatal(err)
}

buf := make([]byte, len(message))
n, err := r.Read(buf)
if err != nil {
	log.Fatal(err)
} else if n != len(message) {
	log.Fatalf("Read %d bytes, expected %d", n, len(message))
}
fmt.Printf("Received: %s\n", buf)

outputBytes, _ := ioutil.ReadAll(stdout)
outputString := string(outputBytes)
fmt.Println("Output: ", outputString)

