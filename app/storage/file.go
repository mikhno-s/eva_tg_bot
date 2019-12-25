package storage

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/zelenin/go-tdlib/client"
)

type DBFile struct {
	File *os.File
	Path string
}

func (file *DBFile) Init() (err error) {
	file.File, err = os.OpenFile(file.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	return
}

func (file *DBFile) Close() error {
	return file.File.Close()
}

func (file *DBFile) ReadMessages() (messages []*client.Message, err error) {
	scanner := bufio.NewScanner(file.File)

	fStat, err := file.File.Stat()
	if err != nil {
		return
	}
	if fStat.Size() > 0 {
		for scanner.Scan() {
			m := client.Message{}
			err = json.Unmarshal(scanner.Bytes(), &m)
			if err != nil {
				return
			}
			messages = append(messages, &m)
		}
	}
	return
}

func (file *DBFile) WriteMessages(messages []*client.Message) (err error) {

	for _, m := range messages {
		marhMessage, err := m.MarshalJSON()
		if err != nil {
			return err
		}
		_, err = file.File.Write(marhMessage)
		if err != nil {
			return err
		}
		_, err = file.File.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return
}
