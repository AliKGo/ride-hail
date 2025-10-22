package config

import (
	"encoding/json"
	"fmt"
)

func (cfg *Config) printConfig() {
	fmt.Println("-------------------- Config --------------------")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Println("error marshaling config:", err)
		return
	}
	fmt.Println(string(data))
	fmt.Println("------------------------------------------------")
}
