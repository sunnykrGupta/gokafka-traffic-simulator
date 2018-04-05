package main

import (
	"os"
	"fmt"
	"strconv"
	"github.com/Pallinder/go-randomdata"
	"encoding/json"
)

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
	//fmt.Printf("%T\n%q", jsn, jsn)

	if err != nil {
		return nil, err
	}

	return jsonBytes, nil

}


func main(){

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <Total_Profiles>\n",
			os.Args[0])
		os.Exit(1)
	}

	Total_Profiles, _ := strconv.Atoi(os.Args[1])

	for i := 0; i < Total_Profiles; i++{

		jsonBytes, e := profileProvider()

		if e == nil{
			//fmt.Printf("%% %s\n", string(jsonBytes)[:50]) //string(e.Value)[:50]
			fmt.Printf("%s\n", jsonBytes)
		}

		/*if ( i % 10000 == 0){
			fmt.Println(i)
			fmt.Printf("%% %s\n", string(jsonBytes)[:50])
		}*/
	}
}
