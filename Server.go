package main

import (
    "fmt"
    "net"
    "strconv"
)

func server(port int) {
    fmt.Println("Server startup!")
    stringPort := strconv.Itoa(port)
    // start listening to a part
    listener, err := net.Listen("tcp", ":" + stringPort)
    if (err != nil) {
        fmt.Println("Couldn't listen on port " + stringPort)
        return
    }
    for {
        connection, err := listener.Accept()
        fmt.Println("Connection established!")

        if (err != nil) {
            fmt.Println("Failed to establish connection.")
        } else {
            go handleClient(connection)
        }
    }
}

func handleClient(client net.Conn) {
    for {
        message := make([]byte, 10)
        bytes, err := client.Read(message)
        stringMessage := string(message[:bytes])

        if (err != nil) {
            fmt.Println("Error when reading from socket.", err)
            client.Close()
            break
        } else {
            fmt.Print(stringMessage)
            _, errWrite := client.Write([]byte("Message recieved!"))
            if (errWrite != nil) { fmt.Println(errWrite) }
        }
        if (stringMessage == "stop\n") {
            break
        }
        if (stringMessage == "close\n") {
            client.Close()
        }
    }
}

func main() {
    var port int
    _, err := fmt.Scanf("%d", &port) 
    if (err != nil) {
        fmt.Println("Couldn't read user argument.")
    } else {
        server(port)
    }
}
