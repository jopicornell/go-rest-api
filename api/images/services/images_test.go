package services

import (
	"github.com/jopicornell/go-rest-api/pkg/config"
	"github.com/jopicornell/go-rest-api/pkg/servertesting"
	"reflect"
	"testing"
)

func TestNewImageService(t *testing.T) {
	mockedServer := servertesting.Initialize(&config.Config{})
	picturesService := NewImagesService(mockedServer)
	if picturesService == nil {
		t.Errorf("New service should not be null")
	}
	reflectServiceDb := reflect.ValueOf(picturesService).Elem().FieldByName("db")
	if reflectServiceDb.IsNil() {
		t.Errorf("db field in service should not be nil")
	}
}
