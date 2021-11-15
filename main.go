//Alex Zicaro 11/14/2021
//----------------------
//TO RUN
//Navigate to file location, run go build main.go
//Then run main
//You may have to unstall github.com/gorilla/websocket with go get -v github.com/gorilla/websocket

//----------------------
//Once running, navigate to http://127.0.0.1:9000/. To adjust headers I used ModHeader extension on Chrome
//The headers I am looking for are:
//  Action: sub
//  Action: pub
//  Message: ~insert message~

//If these headers are missing, I give the user an error

package main
 
import (
    "fmt"
    "net/http"
    "github.com/gorilla/websocket"
    // "log"
)
//Store the IPs of the subscribers
var subs = make([]string, 0)
//Store websocket of subscribers
var cons []*websocket.Conn = nil

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

// create a handler struct
type HttpHandler struct{}

// Handle subscribe requesrs
func (h HttpHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    // create response binary data
    reqIP := GetIP(req)
    fmt.Printf("Handling request from " + reqIP + "\n")

    //This is how you upgrate the html connection to a websocket connection. Not too sure what that means
    // conn, err := upgrader.Upgrade(res, req, nil)
    // if err != nil {
    //     log.Println(err)
    //     return
    // }
    var data []byte

    //These messages would come through via a JSON sent by the websocket. Depending on what kind of action you secified, you could [sub]scribe or [pub]lish
    // action, p, err := conn.ReadMessage()
    // fmt.Printf("Action: " + string(action) + " P: " + string(p) + "\n")

    // if err != nil {
    //     fmt.Printf(reqIP + " has an error. Do nothing\n")
    //     data = []byte("Hello " + reqIP + "! There is an error with your request. No action was taken. Message: " + string(action) + "\n")
    // } else {
    //     switch string(action){
    //     case "sub":
    //         SubscribeClient(conn, reqIP)
    //         break
    //     case "pub":
    //         PublishToClients(conn, reqIP)
    //         break
    //     default:
    //         fmt.Printf(reqIP + " has an error. Do nothing\n")
    //         data = []byte("Hello " + reqIP + "! You submitted an unknown message. No action was taken. Message: " + string(action) + "\n")
    //     }
    
    // }

    //I didn't fully understand the websocket way(which is the correct way to do this) so I found some other differentator for specifying what action you wanted, headers
    //Loop through all headers and find "Action" and switch based on that
    action := req.Header["Action"]
    if action != nil {
        switch action[0] {
        case "sub": 
            SubscribeClient(nil, reqIP, res)
            break
        case "pub":
            if len(req.Header["Message"]) !=0{
                PublishToClients(req.Header["Message"][0], res)
            } else{
                data = []byte("Hello " + reqIP + "! You do not have a message header for your publish. No action was taken.\n")
            }
        default:
            fmt.Printf(reqIP + " has an error. Do nothing\n")
            data = []byte("Hello " + reqIP + "! You have specified an unknown action. No action was taken. Action: " + action[0] + "\n")
        }
    } else {
        fmt.Printf(reqIP + " has no Action header. Do nothing\n")
        data = []byte("Hello " + reqIP + "! You do not have an Action header\n")
    }    //If the action header was not found, show the client an error
    
    // write `data` to response
    res.Write(data)
}

//Handle subscribe actions
func SubscribeClient(conn *websocket.Conn, reqIP string, res http.ResponseWriter){
    fmt.Printf("Handling request from " + reqIP + "\n")
    var data []byte
    //Check if this client has already subscribed
    if !contains(subs, reqIP) {
        fmt.Printf("Appending IP to sub list %v \n", subs)
        subs = append(subs, reqIP)
        // cons = append(cons, conn)
        fmt.Printf("%v\n", subs)
        data = []byte("Hello " + reqIP + "! You have been subscribed.\n")
    } else {
        fmt.Printf(reqIP + " is already subscribed")
        data = []byte("Hello " + reqIP + "! You are already subscribed.\n")
    }
    // conn.WriteMessage(websocket.TextMessage, data)
    res.Write(data)
}

func PublishToClients(message string, res http.ResponseWriter){
    //I have the IP addresses to the clients but I do not have a websocket connection. Tried to go this route with the commented code but was not sure how to initiate the connection from the client
    //I saw some examples online of mkaing some kind of websocket connection using a static html page but did not fully understand it
    //Right now I just write the message from the header to the webpage. Not much but it's somethin
    
    res.Write([]byte(message))
    // for _, con := range cons {
    //     con.WriteMessage(1, []byte(message))
    // }
}


// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
    for _, v := range s {
        if v == str {
            return true
        }
    }
    return false
}

// GetIP gets a requests IP address by reading off the forwarded-for
func GetIP(r *http.Request) string {
    forwarded := r.Header.Get("X-FORWARDED-FOR")
    if forwarded != "" {
        return forwarded
    }
    return r.RemoteAddr
}
 
func main() {
    fmt.Printf("Creating server on http://127.0.0.1:9000\n")
    // create handler
    handler := HttpHandler{}
    // http.HandleFunc("/ws", handler)

    // listen & serve
    http.ListenAndServe(":9000", handler)
}

