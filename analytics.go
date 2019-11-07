package main

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"os"
	"time"
)

var influxdb client.Client

func track(method string) {
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database: os.Getenv("INFLUXDB_NAME"),
	})
	point, err := client.NewPoint("smmsbot", map[string]string{
		"method": method,
	}, map[string]interface{}{
		"c": 1,
	}, time.Now())
	if err != nil {
		return
	}
	bp.AddPoint(point)
	influxdb.Write(bp)
}

func init() {
	var err error
	influxdb, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     os.Getenv("INFLUXDB_ADDR"),
		Username: os.Getenv("INFLUXDB_USER"),
		Password: os.Getenv("INFLUXDB_PASS"),
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
	}
}
