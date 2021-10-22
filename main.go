package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

func main() {
	creds := credentials{}
	creds.loadFromEnv()

	auth := smtp.PlainAuth("", creds.mailUser, creds.mailPass, creds.mailServer)

	http.HandleFunc("/send", func(rw http.ResponseWriter, r *http.Request) {
		var body *message
		var err error
		err = json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			log.Println("ERROR - could not parse the message")
			rw.WriteHeader(500)
			return
		}
		err = smtp.SendMail(creds.mailServer, auth, creds.mailUser, creds.mailRecipients, body.prepare())
		if err != nil {
			log.Println("ERROR - could not send the message")
			rw.WriteHeader(500)
			return
		}
	})
	address := os.Getenv("NS_ADDR")
	port := os.Getenv("NS_PORT")
	if address == "" || port == "" {
		panic("pass server address and port through environment")
	}
	log.Fatalln(http.ListenAndServe(address+":"+port, nil))

}

type credentials struct {
	mailServer     string
	mailUser       string
	mailPass       string
	mailRecipients []string
}

func (c *credentials) loadFromEnv() {
	c.mailServer = os.Getenv("NS_MAIL_SERVER")
	c.mailUser = os.Getenv("NS_MAIL_USER")
	c.mailPass = os.Getenv("NS_MAIL_PASS")
	c.mailRecipients = append(c.mailRecipients, os.Getenv("NS_DEFAULT_RECIPIENT"))

	if c.mailServer == "" || c.mailUser == "" || c.mailPass == "" || c.mailRecipients[0] == "" {
		panic("pass credentials through environment")
	}
}

type message struct {
	Subject     string `json:"subject"`
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
}

func (m *message) prepare() []byte {
	var preparedMessage []byte
	var preparedSubject string
	var preparedContent string

	preparedSubject = fmt.Sprintf("Subject: %s\n", m.Subject)
	preparedContent = fmt.Sprintf("%s says:\n\n%s", m.ServiceName, m.Message)

	preparedMessage = []byte(preparedSubject + preparedContent)
	return preparedMessage
}
