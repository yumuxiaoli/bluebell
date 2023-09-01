package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
)

var sf *snowflake.Node

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	sf, err = snowflake.NewNode(machineID)
	return
}

func GenID() int64 {
	return sf.Generate().Int64()
}

func main() {
	if err := Init("2023-08-21", 1); err != nil {
		fmt.Printf("init failed err:%v\n", err)
		return
	}
	id := GenID()
	fmt.Println(id)
}
