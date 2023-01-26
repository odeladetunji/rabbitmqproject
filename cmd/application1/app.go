package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
)

func main() {
	router := gin.Default();

	router.POST("/message", func(c *gin.Context) {
		var errMessage error = connectToRabbitMQ(c);
		if errMessage != nil {
			c.JSON(400, gin.H{
				"status": http.StatusBadRequest,
				"time": time.Now(),
			});

			return;
		}

		c.JSON(200, "");
	});

	if err := router.Run(":8092"); err != nil {
		fmt.Println(err.Error())
	}
}

func connectToRabbitMQ(c *gin.Context) error {

	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err.Error())
		return errors.New(err.Error());
	}

	defer connection.Close();

	//Open A channel
	connectionChannel, errCh := connection.Channel();
	if errCh != nil {
		return errors.New(errCh.Error())
	}

	defer connectionChannel.Close();

	quee, errQuee := connectionChannel.QueueDeclare("Messaging-Quee", false, false, false, false, nil)
	if errQuee != nil {
		return errors.New(errQuee.Error())
	}

	type MessageDto struct {
		Sender string `json:"sender"`
		Reciever string `json:"receiver:"`
		Message string `json:"message"`
	}

	var messageToSend MessageDto;
	if errEn := c.BindJSON(&messageToSend); errEn != nil {
		return errors.New(errEn.Error());
	}
	
	if len(messageToSend.Message) == 0 {
		return errors.New("Message is required")
	}

	if len(messageToSend.Reciever) == 0 {
		return errors.New("Reciever is required")
	}

	if len(messageToSend.Sender) == 0 {
		return errors.New("Sender is required")
	}

	jsonMessage, errJson := JSON.Marshal(messageToSend);
	if errJson != nil {
		return errors.New(errJson.Error());
	}

	errChan = connectionChannel.Publish(
		"",
		"Messaging-Quee",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(jsonMessage),
		},
	)

	if errChan != nil {
		return errors.New(errChan.Error());
	}

    return nil;
}



