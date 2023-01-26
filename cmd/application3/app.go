package main

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"fmt"
)

func main() {
	router := gin.Default();

	router.GET("/message/list", func(c *gin.Context) {
		var listOfMessages error = readFromRedis(c);
		if errMessage != nil {
			c.JSON(400, gin.H{
				"status": http.StatusBadRequest,
				"time": time.Now(),
			});

			return;
		}

		c.JSON(200, listOfMessages);
	});

	if err := router.Run(":8092"); err != nil {
		fmt.Println(err.Error())
	}
}

func readFromRedis(c *gin.Context){
	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	});

	val, err := rdb.Get(ctx, "key").Result();
    if err != nil {
        return nil, errors.New(err.Error())
    }

	return val, nil;

}


