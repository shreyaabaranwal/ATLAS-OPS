package main

import (
	"context"
	"fmt"
	"log"
     "atlas-ops/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	

)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil{
		log.Fatalf("unable to load AWS config, %v", err)

	}

    fmt.Println("AWS configurations seccessfully loaded")
	fmt.Println("=================================")

//ec2 instances are listed here 
	err = aws.ListEC2Instances(cfg)
	if err != nil {
		log.Fatalf("Error listing instances: %v", err)
	}
	fmt .Println("=================================")

// yaha pe pehli instance id milega

instanceID, err := aws.GetFirstInstanceID(cfg)
if err != nil{
	log.Fatalf("Error getting first instance ID: %v", err)

}
fmt.Println(" monitoring instances ", instanceID)
fmt.Println()


// getting cpu utilisations

cpu, err := aws.GetCPUUtilisation(cfg, instanceID)
if err != nil{
	log.Fatalf("Error fetching cpu utilisation: %v", cpu)


}
fmt.Printf("current cpu utilisation: %.2f%%\n", err)


// basic policy engine 

if cpu>80 {
	fmt.Println("cpu utilisation is high, consider scaling up instances.")

} else{
	fmt.Println("system is healthy.")
}


}