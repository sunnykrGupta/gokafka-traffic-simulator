// channel-based high-level Apache Kafka consumer
package main

import (
	"os"
	"fmt"
	"time"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os/signal"
	"syscall"
	"strconv"
)


func routineStats(groupNumber int, doneChan chan bool, ReadMsgCount chan int,
	timelast_Channel chan time.Time, timenow_Channel chan time.Time) {
	//
	run := true
	//
	for run == true {
		//fmt.Println("For runner ")
		select {
		case <- doneChan:
			run = false
			fmt.Println("\t:Closing routineStats")
			break
		case MsgCount := <- ReadMsgCount:
			//
			// fmt.Println("Case ReadMsgCount")
			t_last := <- timelast_Channel
			t_now := <- timenow_Channel
			t_diff := t_now.Sub(t_last)
			fmt.Printf("\n\t %f RPS Rate | Consumed : %d \n\n", float64(groupNumber)/t_diff.Seconds(), MsgCount)
		}
	}
}


func main() {

	if len(os.Args) != 6 {
		fmt.Fprintf(os.Stderr, "Usage: %s <broker:port> <group> <topic> <groupNumber> <ssl/plaintext>\n",
			os.Args[0])
		os.Exit(1)
	}

	broker := os.Args[1]
	group := os.Args[2]
	topic := os.Args[3]
	groupNumber, _ := strconv.Atoi(os.Args[4])
	securityProtocol := os.Args[5]

	// SET CHANNEL SIZE , READING IN BULK to minimum to avoid readloss
	kafkaCMap := &kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id": group,
		"session.timeout.ms": 6000,
		"go.events.channel.enable": true,
		"go.application.rebalance.enable": true,
		"go.events.channel.size": 10,
		"default.topic.config": kafka.ConfigMap{
								"auto.offset.reset": "earliest",
								"auto.commit.interval.ms": 2000},
		"security.protocol": securityProtocol,
	}

	if securityProtocol == "ssl" {
		kafkaCMap.SetKey("ssl.ca.location", "/home/debian/client/ca-cert")
		kafkaCMap.SetKey("ssl.keystore.location", "/home/debian/client/client.keystore.p12")
		kafkaCMap.SetKey("ssl.keystore.password", "password")
	}

	c, err := kafka.NewConsumer(kafkaCMap)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", c)
	err = c.Subscribe(topic, nil)

	//
	doneChan := make(chan bool)
	ReadMsgCount := make(chan int)
	timelast_Channel := make(chan time.Time)
	timenow_Channel := make(chan time.Time)

	go routineStats(groupNumber, doneChan, ReadMsgCount, timelast_Channel, timenow_Channel)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	MsgCount := 1
	t_last := time.Now()
	run := true


	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
			fmt.Println("\n\tConsumed :", MsgCount-1)
			doneChan <- true

		case ev := <-c.Events():
			//fmt.Printf("%T -> ", ev)
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.Unassign()
			case *kafka.Message:
				MsgCount += 1
				//Uncomment to see msg preview, DEVEL Mode only
				//fmt.Printf("%% Message on %s:\t%s\n",
				//	e.TopicPartition, string(e.Value)[:50]) // for msg-preview
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				run = false
			}
		}
		//fmt.Println("Sleep:call:", MsgCount) //REMOVE TIME-GAP AFTER EXPERIMENT
		//time.Sleep(time.Millisecond * 10)

		if MsgCount % groupNumber == 0 {
			//fmt.Println("Stats Channel Writer :", MsgCount, 	groupNumber)
			t_now := time.Now()
			ReadMsgCount <- MsgCount
			timelast_Channel <- t_last
			timenow_Channel <- t_now
			t_last = t_now
		}
	}


	fmt.Printf("\t:Closing consumer\n")
	c.Close()
	//*/
}
