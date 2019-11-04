package database

import "testing"

func TestMySQL_GetDB(t *testing.T) {
	t.Run("returning db when it is already initialized", testDBWhenInitialized)
	t.Run("returning new connection when db not initialized", testDBWhenNotInitialized)
}

func testDBWhenInitialized(t *testing.T) {

}

func testDBWhenNotInitialized(t *testing.T) {

}
