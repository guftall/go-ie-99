package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "0.0.0.0", "http service address")
var port = os.Getenv("PORT")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		var command message
		log.Printf("recv: %s", msg)
		err = json.Unmarshal([]byte(msg), &command)

		if err != nil {
			log.Print("unmarshal message failed", err)
		}
		result := runCommand(command)

		log.Print("command result: ", result)
		err = c.WriteMessage(mt, []byte(result))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func runCommand(msg message) string {

	switch msg.Action {
	case "echo":
		{
			return msg.Action
		}
	case "count":
		{
			count := countRecords()
			res := strconv.Itoa(count)
			return res
		}
	case "read_identifier":
		{
			var key string
			for _, val := range msg.Parameters {
				if val.Key == "key" {
					key = val.Value
					break
				}
			}
			identifier := readIdentifier(key)
			return identifier
		}
	case "is_identifier_exist":
		{
			var identifier string
			for _, val := range msg.Parameters {
				if val.Key == "identifier" {
					identifier = val.Value
					break
				}
			}
			exist := isIdentifierExist(identifier)
			if exist {
				return "yes"
			}

			return "no"
		}
	case "insert":
		{
			var key string
			var identifier string
			for _, val := range msg.Parameters {
				if val.Key == "identifier" {
					identifier = val.Value
				}
				if val.Key == "key" {
					key = val.Value
				}
			}

			exist := isIdentifierExist(identifier)
			if exist {
				return "already_exist"
			}

			insertPublicKey(key, identifier)

			return "inserted"
		}
	default:
		{
			return "unrecognized command"
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	wss := os.Getenv("websocketschema")
	homeTemplate.Execute(w, wss+"://"+r.Host+"/websocket")
}

func initializeServer() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/websocket", echo)
	http.HandleFunc("/", home)
	lisAddr := *addr + ":" + port
	log.Print("listening on ", lisAddr)
	log.Fatal(http.ListenAndServe(lisAddr, nil))
}

type message struct {
	Action     string      `json:"action"`
	Parameters []parameter `json:"parameters"`
}

type parameter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
    };
    
    function open(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("open").onclick = open
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
    open()
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
