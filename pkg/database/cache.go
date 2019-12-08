package database

type Cache interface {
	InitializeClient()
	GetStruct(key string, ifc interface{})
	SetStruct(key string, m interface{})
}
