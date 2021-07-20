package main

import (
	"testing"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
)

func TestOutputKafkaRAW(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer := mocks.NewAsyncProducer(t, config)
	producer.ExpectInputAndSucceed()

	output := NewKafkaOutput("", &OutputKafkaConfig{
		producer: producer,
		Topic:    "test",
		UseJSON:  false,
	}, nil)

	output.PluginWrite(&Message{Meta: []byte("1 2 3\n"), Data: []byte("GET / HTTP1.1\r\nHeader: 1\r\n\r\n")})

	resp := <-producer.Successes()

	data, _ := resp.Value.Encode()

	if string(data) != "1 2 3\nGET / HTTP1.1\r\nHeader: 1\r\n\r\n" {
		t.Errorf("Message not properly encoded: %q", data)
	}
}

func TestOutputKafkaJSON(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer := mocks.NewAsyncProducer(t, config)
	producer.ExpectInputAndSucceed()

	output := NewKafkaOutput("", &OutputKafkaConfig{
		producer: producer,
		Topic:    "test",
		UseJSON:  true,
	}, nil)

	output.PluginWrite(&Message{Meta: []byte("1 2 3\n"), Data: []byte("GET / HTTP1.1\r\nHeader: 1\r\n\r\n")})

	resp := <-producer.Successes()

	data, _ := resp.Value.Encode()

	if string(data) != `{"Req_URL":"","Req_Type":"1","Req_ID":"2","Req_Ts":"3","Req_Method":"GET"}` {
		t.Error("Message not properly encoded: ", string(data))
	}
}

func TestOutputKafkaResponseWithReasonPhraseJSON(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer := mocks.NewAsyncProducer(t, config)
	producer.ExpectInputAndSucceed()

	output := NewKafkaOutput("", &OutputKafkaConfig{
		producer: producer,
		Topic:    "test",
		UseJSON:  true,
	}, nil)

	output.PluginWrite(&Message{Meta: []byte("2 2 3\n"), Data: []byte("HTTP/1.1 200 OK\r\n")})

	resp := <-producer.Successes()

	data, _ := resp.Value.Encode()

	if string(data) != `{"Req_URL":"200","Req_Type":"2","Req_ID":"2","Req_Ts":"3","Req_Method":"HTTP/1.1"}` {
		t.Error("Message not properly encoded: ", string(data))
	}
}


func TestOutputKafkaResponseWithoutReasonPhraseJSON(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer := mocks.NewAsyncProducer(t, config)
	producer.ExpectInputAndSucceed()

	output := NewKafkaOutput("", &OutputKafkaConfig{
		producer: producer,
		Topic:    "test",
		UseJSON:  true,
	}, nil)

	output.PluginWrite(&Message{Meta: []byte("2 3 4\n"), Data: []byte("HTTP/1.1 404\r\n")})

	resp := <-producer.Successes()

	data, _ := resp.Value.Encode()

	if string(data) != `{"Req_URL":"404","Req_Type":"2","Req_ID":"3","Req_Ts":"4","Req_Method":"HTTP/1.1"}` {
		t.Error("Message not properly encoded: ", string(data))
	}
}
