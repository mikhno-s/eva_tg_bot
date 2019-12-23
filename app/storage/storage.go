package storage

import (
	"github.com/zelenin/go-tdlib/client"
)

type Row struct {
	Row []byte
}

type MessageRow interface {
	GetRows() []*client.Message
	GetRow(n int) *client.Message
}

type DBFile struct {
	Path    string
	Content []Row
}

// read

// write
