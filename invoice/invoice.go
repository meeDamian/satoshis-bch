package invoice

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

type (
	Pixel struct {
		Coords [2]string `json:"coordinates"`
		Color  string    `json:"color"`
	}

	Picture []Pixel

	NewOrderMsg [2]interface{}
	NewOrderResult [2]json.RawMessage
)

var (
	allowedColors = []string{
		"#ffffff",
		"#e4e4e4",
		"#888888",
		"#222222",
		"#e4b4ca",
		"#d4361e",
		"#db993e",
		"#8e705d",
		"#e6d84e",
		"#a3dc67",
		"#4aba38",
		"#7fcbd0",
		"#5880a8",
		"#3919d1",
		"#c27ad0",
		"#742671",
	}
)

func GetInvoice(pic Picture) (invoice string, err error) {
	c, _, err := websocket.DefaultDialer.Dial("wss://api.satoshis.place/socket.io/?EIO=3&transport=websocket", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	c.EnableWriteCompression(true)

	defer c.Close()

	bytes, err := json.Marshal(NewOrderMsg{"NEW_ORDER", pic})
	if err != nil {
		log.Println("marshall err:", err)
		return "", err
	}

	err = c.WriteMessage(websocket.TextMessage, append([]byte("42"), bytes...))
	if err != nil {
		log.Println("write2:", err)
		return
	}

	log.Println("noerr")

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return "", err
		}

		if !strings.HasPrefix(string(message), "42") {
			continue
		}

		log.Println("has 42")
		message = message[2:]

		log.Printf("recv: %s", message)

		var v NewOrderResult
		err = json.Unmarshal(message, &v)
		if err != nil {
			log.Println("json1", err)
			return "", err
		}

		log.Println(string(v[0]), "NEW_ORDER_RESULT")

		if string(v[0]) == `"NEW_ORDER_RESULT"` {
			log.Println("in")
			type bla struct {
				Data struct {
					Invoice string `json:"paymentRequest"`
				} `json:"data"`
			}

			var b bla
			err = json.Unmarshal(v[1], &b)
			if err != nil {
				log.Println("json2", err)
				return "", err
			}

			invoice = b.Data.Invoice

			break
		}
	}

	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return
	}

	log.Println("fin")

	return
}
