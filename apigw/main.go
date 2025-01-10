package main

import (
	"apigw/routing"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	routing.Routing()
	fmt.Println("API Gateway is running on port", viper.GetString("PORT"))
	log.Fatal(http.ListenAndServe(":"+viper.GetString("PORT"), nil))
}
