// function-based Apache Kafka producer
package main

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"syscall"
	"os/signal"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/Pallinder/go-randomdata"
	"encoding/json"
)


// define const of this producer

//----------------

//struct to keep MSG
type ProfilePara struct {
    Bio *randomdata.Profile
    P1 string
	P2 string
	P3 string
	P4 string
	P5 string
}

func profileProvider() ([]byte, error) {

	profilePara := ProfilePara{
			Bio: randomdata.GenerateProfile(randomdata.RandomGender),
			P1: randomdata.Paragraph() + randomdata.Paragraph() + randomdata.Paragraph() + randomdata.Paragraph(),
			P2: randomdata.Paragraph() + randomdata.Paragraph() + randomdata.Paragraph() + randomdata.Paragraph(),
			P3: randomdata.Paragraph() + randomdata.Paragraph() + randomdata.Paragraph() + randomdata.Paragraph(),
			P4: randomdata.Paragraph() + randomdata.Paragraph() + randomdata.Paragraph() + randomdata.Paragraph(),
			P5: randomdata.Paragraph() + randomdata.Paragraph() + randomdata.Paragraph() + randomdata.Paragraph(),
	}
	jsonBytes, err := json.Marshal(profilePara)

	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

//----------------

func routineStats(msgBurst int, statsDone chan bool, counterChannel chan int,
                timelast_Channel chan time.Time, timenow_Channel chan time.Time ){
    //infinite loop for channel read wait
	run := true
    for run == true {

		select {
		case <- statsDone:
			run = false
			fmt.Println("Closed Stats Goroutine")
			break

		case t_last := <- timelast_Channel:
	        t_now := <- timenow_Channel
	        counter := <- counterChannel
            // result will be duration
            t_diff := t_now.Sub(t_last)
            fmt.Println(t_last.Format("15:04:05.999"), t_diff.Seconds()*1000)
            fmt.Println(float64(msgBurst)/t_diff.Seconds(), " RPS Rate | ", "Produced :", counter)
            fmt.Println()
            //
		}
    }
}



func main() {

	if len(os.Args) != 8 {
		fmt.Fprintf(os.Stderr, "Usage: %s <broker:port> <topic> <msgBurst> <total> <ssl/plaintext> <lingerMs> <waitMs>\n",
			os.Args[0])
		os.Exit(1)
	}

	broker := os.Args[1]
	topic := os.Args[2]
	msgBurst, _ := strconv.Atoi(os.Args[3])
	total, _ := strconv.Atoi(os.Args[4])
	securityProtocol := os.Args[5]
	lingerMs, _ := strconv.Atoi(os.Args[6])
	waitMs, _ := strconv.Atoi(os.Args[7])

	//linger.ms : Delay in milliseconds to wait for messages in the producer queue to accumulate before
	//"message.max.bytes": 1000,000, Maximum Kafka protocol request message(batch) size.
	//"queue.buffering.max.kbytes" 1048576, (1GB) Maximum total message size sum allowed on the producer queue.

	kafkaCMap := &kafka.ConfigMap{"bootstrap.servers": broker,
								"request.required.acks": 1,
								"queue.buffering.max.kbytes": 204800,
								"queue.buffering.max.messages": 100000,
								"linger.ms": lingerMs,
								"security.protocol": securityProtocol,
								"compression.codec": "gzip",
	}

	if securityProtocol == "ssl" {
		kafkaCMap.SetKey("ssl.ca.location", "/home/debian/client/ca-cert")
		kafkaCMap.SetKey("ssl.keystore.location", "/home/debian/client/client.keystore.p12")
		kafkaCMap.SetKey("ssl.keystore.password", "password")
	}

	p, err := kafka.NewProducer(kafkaCMap)

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Producer %v\n", p)

	//Channel declaration
	//doneChan := make(chan bool)
	statsDone := make(chan bool)
    counterChannel := make(chan int)
    timelast_Channel := make(chan time.Time)
    timenow_Channel := make(chan time.Time)

	// goroutine invoke which process stats
	go routineStats(msgBurst, statsDone, counterChannel, timelast_Channel, timenow_Channel)

	go func(){
		//defer close(doneChan)

		for e := range p.Events(){
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				}
				//else {
				//	fmt.Printf("Delivery message to topic %s [%d] at offset %v\n",
				//		 *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset )
				//}

				//return	BIGGEST FRAUD LIED DURING THE WHOLE TRIAL.

			default :
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()


	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	msgCnt := 1
	run := true
	t_last := time.Now()

	//To produce same packet repetitively
	jsonBytes, _ := profileProvider()

	for msgCnt <= total {

		select {
			case sig := <- sigchan:
				fmt.Printf("Caught Signal %v: terminating\n", sig)
				fmt.Printf("\n\tProduced : %d\n", msgCnt-1)
				run = false
				//terminate goroutine stats
				statsDone <- true

			default:
				//AVOID randomdata
				//jsonBytes, err := profileProvider()

				err = p.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
					Value:          []byte(jsonBytes),
				}, nil)

				if err == nil {
					msgCnt++
				}else {
					fmt.Println("Produce Err : ", err)
					fmt.Printf("Wait for delivery : %d ms\n\n", waitMs)
					time.Sleep(time.Millisecond * time.Duration(waitMs))
				}

				//Periodically invoke stats crunchers
				if msgCnt % msgBurst == 0 {
		            t_now :=  time.Now()
					//Write to channels
		            timelast_Channel <- t_last
		            timenow_Channel <- t_now
		            counterChannel <- msgCnt
		            t_last = t_now
				}
		}

		if run == false{
			break
		}

		if msgCnt == total {
			statsDone <- true
		}
	}


	n := p.Flush(30000)
	fmt.Println("Outstanding events unflushed : ", n)

	p.Close()
}
