package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)


func init(){
	Load(".env")
}


func Load(file string){
	err := godotenv.Load(file)
	if err != nil {
		log.Println("can't find .env file")
	}
}

func GetEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func GetEnvAsInt(key string, defaultValue int) int {
    valueStr := GetEnv(key, "")
    if value, err := strconv.Atoi(valueStr); err == nil {
        return value
    }
    return defaultValue
}
