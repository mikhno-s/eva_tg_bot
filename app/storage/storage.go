package storage

import (
	"github.com/zelenin/go-tdlib/client"
)

type Storage interface {
	Init() (err error)
	Close() (err error)
	ReadMessages() (messages []*client.Message, err error)
	WriteMessages(messages []*client.Message) (err error)
}

func InitStorage(storage Storage) (Storage, error) {
	err := storage.Init()
	if err != nil {
		return nil, err
	}
	return storage, nil
}
