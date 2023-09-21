package main

import (
	"time"
	"fmt"
	"log"
	"encoding/json"
	"io/ioutil"
	"go-tracking/models"
	_packageClient "go-tracking/client"
	"go-tracking/client/httpsocket"
	_packageUcase "go-tracking/client/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TrackingData struct {
	Id        string  `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.File("/", "web/index.html")

	pc, err := _packageClient.NewRabbitMQClient("amqp://admin:1sampai8@localhost:5672/")
	if err != nil {
		log.Panicln(err)
	}
	defer pc.Close()

	go func() {

		// Read JSON data from "data/tracking.json"
		data, err := ioutil.ReadFile("data/tracking.json")
		if err != nil {
			log.Println("Error reading JSON data:", err)
			return
		}
	
		// Parse JSON data into a slice of TrackingData
		var trackingData []TrackingData
		if err := json.Unmarshal(data, &trackingData); err != nil {
			log.Println("Error unmarshaling JSON data:", err)
			return
		}
		
		// Now, you can loop through the trackingData slice and use the models package as needed.
		for _, item := range trackingData {
			tracking := &models.Tracking{
				Id:        item.Id,
				Latitude:  item.Latitude,
				Longitude: item.Longitude,
			}
			fmt.Println(*tracking)

			time.Sleep(1 * time.Second)
			pc.Publish(*tracking)
	
		}
	}()

	pu := _packageUcase.NewPackageUsecase(pc)

	httpsocket.NewPackageHandler(e, pu)

	e.Logger.Fatal(e.Start(":1323"))
}
