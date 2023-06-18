package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/corluk/believealbumler-elastic-consumer/models"
	"github.com/segmentio/kafka-go"
)

func setAuthorizationHeader(req *http.Request) {

	user := os.Getenv("ELASTIC_USER")
	pass := os.Getenv("ELASTIC_PASS")
	encoded := base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", encoded))
}
func DoRequest(req *http.Request) (*http.Response, error) {

	setAuthorizationHeader(req)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func createReleaseIndex(index string) error {
	uri := os.Getenv("ELASTIC_URI") + "/" + index
	reqIsExists, err := http.NewRequest("HEAD", uri, nil)

	if err != nil {
		return err
	}

	respIsExists, err := DoRequest(reqIsExists)
	if err != nil {
		return err
	}

	if respIsExists.StatusCode != int(200) {
		reqNew, err := http.NewRequest("PUT", uri, nil)
		if err != nil {
			return err
		}
		respNew, err := DoRequest(reqNew)
		if err != nil {
			return err
		}

		if respNew.StatusCode != 200 {

			return errors.New("bad response code " + strconv.Itoa(respNew.StatusCode))
		}

	}

	return nil
}
func saveAudioItem(index string, key string, audioItem models.AudioItem) error {

	uri := os.Getenv("ELASTIC_URI") + "/" + index + "/_doc/" + key

	reqExists, err := http.NewRequest("HEAD", uri, nil)
	if err != nil {
		return err
	}
	resp, err := DoRequest(reqExists)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		// delet e
		reqDelete, err := http.NewRequest("DELETE", uri, nil)
		if err != nil {
			return err
		}
		respDelete, err := DoRequest(reqDelete)
		if err != nil {
			return err
		}
		if respDelete.StatusCode != 200 {
			return errors.New("bad response code " + strconv.Itoa(respDelete.StatusCode))

		}
	}

	b, err := json.Marshal(audioItem)
	if err != nil {
		return err
	}
	uri = os.Getenv("ELASTIC_URI") + "/" + index + "/_create/" + key
	reqNew, err := http.NewRequest("PUT", uri, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	respNew, err := DoRequest(reqNew)
	if err != nil {
		return err
	}
	//str, err := ioutil.ReadAll(respNew.Body)
	/*if err != nil {
		return err
	}
	*/
	//fmt.Println(str)
	if respNew.StatusCode != 201 {
		return errors.New("bad response code " + strconv.Itoa(respNew.StatusCode))

	}

	return nil

}
func run(item models.Message) error {

	return saveAudioItem("believe_albums", item.Key.Hex(), item.Value)

}

func readKafka(channel chan models.AudioItem) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(os.Getenv("KAFKA_BROKERS"), ","),

		Topic: "test-believealbum-search-with-key",
	})
	reader.SetOffset(0)
	fmt.Println("starting thread")
	total := 0
	for {
		var message models.Message

		incomingMessage, err := reader.ReadMessage(context.Background())
		total = total + 1
		date := time.Now()

		fmt.Printf("message time %s\n", date.String())

		if err != nil {
			fmt.Printf("error reading message %s", err.Error())
			continue
		}
		err = json.Unmarshal(incomingMessage.Value, &message)

		if err != nil {
			fmt.Printf("unmarshall %s", err.Error())

			continue
		}
		err = run(message)
		if err != nil {
			panic(err)
		}

	}

}
func main() {

	os.Setenv("KAFKA_BROKERS", "localhost:9092")
	os.Setenv("ELASTIC_USER", "elastic")
	os.Setenv("ELASTIC_PASS", "AeYhdH2BLhSjmWp3sSvq")
	os.Setenv("ELASTIC_URI", "http://localhost:9200")

	err := createReleaseIndex("believe_albums")
	if err != nil {
		panic(err)
	}

	chanAudioItem := make(chan models.AudioItem, 1)
	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, os.Interrupt, syscall.SIGINT)
	go readKafka(chanAudioItem)
	select {

	case audioItem := <-chanAudioItem:
		fmt.Println(audioItem)
	case <-gracefulStop:
		os.Exit(0)

	}
}
