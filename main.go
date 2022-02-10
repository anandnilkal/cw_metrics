package main

import (
	"context"
	"fmt"
)

func main() {

	//cw_client, err := NewCloudwatchClient(context.Background(), "us-west-2")
	cw_client, err := NewCloudwatchClient(context.Background(), "us-west-2")
	if err != nil {
		fmt.Errorf("Unable to create cloudwatch client")
		return
	}

	cw_client.getMetrics(context.Background())
}
