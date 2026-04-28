package storage

import "encoding/json"

type DataType string

const (
	TypeCredential DataType = "credential"
	TypeText       DataType = "text"
	TypeCard       DataType = "card"
	TypeBinary     DataType = "binary"
)

// Data — универсальная модель данных.
type Data struct {
	ID    string          `json:"id"`
	User  string          `json:"user"`
	Type  DataType        `json:"type"`
	Value json.RawMessage `json:"value"`
	Meta  string          `json:"meta"`
}
