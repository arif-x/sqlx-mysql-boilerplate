package uuid

import (
	"log"

	"github.com/google/uuid"
)

func StringToUUID(stringUUID string) (uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(stringUUID)
	if err != nil {
		log.Fatal(err)
	}

	return parsedUUID, err
}
