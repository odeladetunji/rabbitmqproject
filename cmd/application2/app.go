package main

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"fmt"
)

func main() {
	router := gin.Default();

	readRabbitMQMessages();

	if err := router.Run(":8091"); err != nil {
		fmt.Println(err.Error())
	}
}


func readRabbitMQMessages() {

	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err.Error())
	}

	defer connection.Close();

	//Open A channel
	connectionChannel, errCh := connection.Channel();
	if errCh != nil {
		fmt.Println(errCh.Error());
	}

	defer connectionChannel.Close();

	incommingMessages, errQuee := connectionChannel.Consume("Messaging-Quee", "", true, false, false, false, nil)
	if errQuee != nil {
		fmt.Println(errQuee.Error())
	}

	listeningChannel := make(chan bool);
	listenerRoutine := func() {
		for mess := range incommingMessages {
			saveToRedis(mess.Body);
			fmt.Println(mess.Body)
		}
	}();

	go listenerRoutine();

	<-listeningChannel
}

func saveToRedis(byteArray []byte) {
	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	});
	
	err := rdb.Set(ctx, "message", byteArray, 0).Err();
	if err != nil {
		panic(err)
	}

	
}




