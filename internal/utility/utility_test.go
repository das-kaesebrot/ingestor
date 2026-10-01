package utility_test

import (
	"testing"
	stdlibuuid "uuid"

	"dev.kaesebrot.eu/go/ingestor/internal/utility"
)

func TestToGoogleUUID(t *testing.T) {
	id := stdlibuuid.NewV7()
	googleuuid := utility.ToGoogleUUID(id)

	if id.String() != googleuuid.String() {
		t.Fatalf("%s != %s", id.String(), googleuuid.String())
	}
}
