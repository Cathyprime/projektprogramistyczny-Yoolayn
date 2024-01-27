package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	amount := 1
	var strip bool

	for _, v := range os.Args[1:] {
		if v == "-s" || v == "--strip" {
			strip = true
			continue
		}
		am, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		amount = am
	}

	for i := 0; i <	amount; i++ {
		if strip {
			str := primitive.NewObjectID().String()
			str = strings.Replace(str, "ObjectID(\"", "", -1)
			str = strings.Replace(str, "\")", "", -1)
			fmt.Println(str)
		} else {
			fmt.Println(primitive.NewObjectID())
		}
	}
}
