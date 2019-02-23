package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"encoding/json"
)

var gcm cipher.AEAD
var nonce []byte
type Message struct {
	Data []byte `json:"data"`
}

func main() {
	key := []byte("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}	
	gcm, err = cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}
	nonce = make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}
	http.HandleFunc("/encrypt", encrypt)
	http.HandleFunc("/decrypt", decrypt)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func parseRequestBody(w http.ResponseWriter, body io.Reader) ([]byte, error) {
	decoder := json.NewDecoder(body)
	var message Message
	err := decoder.Decode(&message)
	return message.Data, err
}

func writeMessageResponse(w http.ResponseWriter, message Message) (error) {
 	w.Header().Set("Content-Type", "application/json")
	messageJson, err := json.Marshal(message)
 	w.Write(messageJson)
	return err
}

func encrypt(w http.ResponseWriter, r *http.Request) {
	plaintext, err := parseRequestBody(w, r.Body)
	fmt.Println(string(plaintext))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var message Message
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	message.Data = ciphertext
	err = writeMessageResponse(w, message)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}	
}

func decrypt(w http.ResponseWriter, r *http.Request) {
	ciphertext, err := parseRequestBody(w, r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		fmt.Println("ciphertext is smaller than nonce size")
	}
	nonce, ciphertext = ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
	}
	var message Message
	message.Data = plaintext
	err = writeMessageResponse(w, message)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}	
}
