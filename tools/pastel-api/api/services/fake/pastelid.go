package fake

import (
	"crypto/sha256"
	"fmt"

	"github.com/pastelnetwork/gonode/tools/pastel-api/api/services/fake/models"
)

// PastelID represents pastelID.
type PastelID struct {
	id string
}

func (pastel PastelID) sign(text string) *models.PastelIDSign {
	return &models.PastelIDSign{
		Signature: pastel.signature(text),
	}
}

func (pastel PastelID) verify(text, signature string) *models.PastelIDVerify {
	var verification string
	if pastel.signature(text) == signature {
		verification = "OK"

	}
	return &models.PastelIDVerify{
		Verification: verification,
	}
}

func (pastel PastelID) signature(text string) string {
	sum := sha256.Sum256([]byte(pastel.id + " " + text))
	return fmt.Sprintf("%x", sum)
}

func newPastelID(id string) *PastelID {
	return &PastelID{
		id: id,
	}
}
