package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
	// "encoding/json"
)

func main() {
	hostname, _ := os.Hostname()

	server := flag.String("server", "tcp://iot.eclipse.org:1883", "The full URL of the MQTT server to connect to")
	topic := flag.String("topic", "fbrnc/test", "Topic to publish the messages on")
	qos := flag.Int("qos", 0, "The QoS to send the messages at")
	retained := flag.Bool("retained", true, "Are the messages sent with the retained flag")
	clientid := flag.String("clientid", hostname+strconv.Itoa(time.Now().Second()), "A clientid for the connection")
	username := flag.String("username", "", "A username to authenticate to the MQTT server")
	password := flag.String("password", "", "Password to match username")
	flag.Parse()

	connOpts := MQTT.NewClientOptions().AddBroker(*server).SetClientID(*clientid).SetCleanSession(true)
	if *username != "" {
		connOpts.SetUsername(*username)
		if *password != "" {
			connOpts.SetPassword(*password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return
	}
	fmt.Printf("Connected to %s\n", *server)

	fmt.Println("Starting server on port 5001")

	http.ListenAndServe(":5001", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()
		var buf bytes.Buffer
		/*
		if err := json.Indent(&buf, b, " >", "  "); err != nil {
			panic(err)
		}
		*/
		fmt.Println(buf.String())
		fmt.Println(b)
		client.Publish(*topic, byte(*qos), *retained, b)
	}))

}
