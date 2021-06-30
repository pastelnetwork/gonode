package artworksearch

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
)

func inIntRange(val int, min *int, max *int) bool {
	if min != nil && val < *min {
		return false
	}
	if max != nil && val > *max {
		return false
	}

	return true
}

func fromBase64(encoded string, to interface{}) error {
	bytes, err := b64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return fmt.Errorf("decode: %s", err)
	}

	if err := json.Unmarshal(bytes, to); err != nil {
		return fmt.Errorf("unmarshal: %s", err)
	}

	return nil
}
