
## Repo contains
- randomdata-producer-from-file #to generate file with random profiles
- cloud
    - consumer_channel_cloud.go     # to consume topic and display reading stats
    - producer_function_cloud.go   #to produce single payload repetitively and display production stats


#### Producer 500,000 records (Producer Client)

> Usage: ./producer_function_cloud <broker:port> <topic> <msgBurst> <total> <ssl/plaintext> <lingerMs> <waitMs>

- broker:port : all Kafka Broker details
- topic : Topic where you want to producer messages
- msgBurst : To display stats after producing these many messages
- total : Total message to produce
- ssl/plaintext : Choose protocol to communicate.
- lingerMs : Time to burst producer messages
- waitMs : Wait time when internal queue is full.

```
$ ./producer_function_cloud broker1:9093,broker2:9093,broker3:9093 MaveRickInfo 10000 500000 ssl 100 2000
```


#### Consume message indefinitely (Consumer Client)

> Usage: ./consumer_channel_cloud <broker:port> <group> <topic> <groupNumber> <ssl/plaintext>

- broker:port : all Kafka Broker details
- group : Consumer Group Name
- topic : Topic where you want to producer messages
- groupNumber : To display stats after reading these many messages
- ssl/plaintext : Choose protocol to communicate. If using 'ssl', configure ssl parameters in code. Refer Kafka Documentation to generate certificate.

```
$ ./consumer_channel_cloud broker1:9093,broker2:9093,broker3:9093 CGrp_01 MaveRickInfo 10000 plaintext
```


-------------------

### Dependency Installation on Linux (Centos Server)

```
sudo yum install epel* -y
sudo yum install -y vim wget telnet net-tools git gcc gcc-c++
sudo yum install -y openssl openssl-devel nload zlib1g-dev zlib1g

sudo yum groupinstall 'Development Tools'
```

#### Install librdkafka to run go-kafka code

[https://github.com/edenhill/librdkafka/]()

```
git clone https://github.com/edenhill/librdkafka.git
cd librdkafka
./configure --prefix /usr
make
sudo make install
```

#### Install gokafka
[https://github.com/confluentinc/confluent-kafka-go]()

#### Install go-randomdata for fake profile generator
[https://github.com/Pallinder/go-randomdata]()



##### For  SSL based support
[https://github.com/edenhill/librdkafka/wiki/Using-SSL-with-librdkafka]()

and Generate Key-pairs for each client and CA pub key (one the broker certificate has been signed with) when to encrpt the data.


---------------


#### Troubleshoot

./producer_function_cloud: error while loading shared libraries: librdkafka.so.1: cannot open shared object file: No such file or directory

```
$ whereis librdkafka
```
Add that directory (/usr/lib) to ld.so.conf file

```
sudo vim /etc/ld.so.conf
sudo ldconfig
```
