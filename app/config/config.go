package config

import (
	"os"
	"strconv"
)

type Config struct {
	ApiID              int32
	ApiHash            string
	ChanName           string
	MessageStoragePath string
	OutputPath         string
}

func InitConfig() *Config {
	// TODO use Cobra
	strAppID, _ := strconv.Atoi(os.Getenv("API_ID"))
	apiID := int32(strAppID)
	apiHash := os.Getenv("API_HASH")
	publicChannelUsername := os.Getenv("CHAN_NAME")
	storageFilePath := os.Getenv("MESSAGE_STORAGE_FILE")
	outputPath := os.Getenv("OUTPUT_FILE")
	config := &Config{
		ApiID:              apiID,
		ApiHash:            apiHash,
		ChanName:           publicChannelUsername,
		MessageStoragePath: storageFilePath,
		OutputPath:         outputPath,
	}
	return config
}
