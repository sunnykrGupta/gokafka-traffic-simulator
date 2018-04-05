
## Folder contains
- utils
  - faker.go #to generate file with random profiles
- producer_cloud_fRead.go   #to produce traffic reading random profiles from file
- producer_cloud_fRead_noGzip.go   #to produce non-compressed traffic reading random profiles from file


#### Configure file (which contains random profiles, we are using file to read profiles before hand prepared payload to avoid slowing down our producer) location inside code before running and build.
>> file, err := os.Open("/home/debian/megafile")

### To pump 500,000 records from file containing 500,0000 random profiles,
>> Usage: ./producer_cloud_fRead <broker:port> <topic> <msgBurst> <total> <ssl/plaintext> <lingerMs> <waitMs>
>> ./producer_cloud_fRead broker1:9093,broker2:9093,broker3:9093  TopicName 10000 500000 ssl 100 5000
