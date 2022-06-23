package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Configuration ...
type configuration struct {
	Host    string `mapstructure:"HOST"`
	Port    int    `mapstructure:"PORT"`
	Debug   bool   `mapstructure:"DEBUG"`
	ConnStr string `mapstructure:"CONNECTION_STRING"`
	Name    string `mapstructure:"SERVICE_NAME"`
}

var (
	conf = &configuration{
		Host: "",
		Port: 8080,
	}
)

func printConf(w http.ResponseWriter, req *http.Request) {
	rawJson, err := json.Marshal(conf)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(rawJson)
}

func main() {
	if err := loadConfig(conf); err != nil {
		log.Fatalf(`Failed to load application configuration: %v`, err)
	} else {
		log.Printf("NAME: '%v'\n", conf.Name)
		log.Printf("HOST: '%v'\n", conf.Host)
		log.Printf("PORT: '%v'\n", conf.Port)
		log.Printf("DEBUG: '%v'\n", conf.Debug)
		log.Printf("CONNECTION_STRING: '%v'\n", conf.ConnStr)
	}

	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	log.Printf("Starting server on: '%s'\n", addr)
	http.HandleFunc("/", printConf)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func loadConfig(conf *configuration) error {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Workaround for the viper to use environment variables without reading a config file
	// Viper #761: Unmarshal non-bound environment variables
	// https://github.com/spf13/viper/issues/761
	envKeysMap := &map[string]interface{}{}
	if err := mapstructure.Decode(conf, &envKeysMap); err != nil {
		return err
	}
	for k := range *envKeysMap {
		if err := viper.BindEnv(k); err != nil {
			return err
		}
	}
	// END (Workaround)

	if err := viper.Unmarshal(conf); err != nil {
		return err
	}

	validate := validator.New()
	return validate.Struct(conf)
}
