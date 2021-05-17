package wdm

import (
	"github.com/pastelnetwork/gonode/dupe-detection/wdm/wrapper"
)

// Wdm compute correlation for specified method. List of supported methods:
//   - `"pearson"`, `"prho"`, `"cor"`: Pearson correlation
//   - `"spearman"`, `"srho"`, `"rho"`: Spearman's
//   - `"kendall"`, `"ktau"`, `"tau"`: Kendall's tau
//   - `"blomqvist"`, `"bbeta"`, `"beta"`: Blomqvist's beta
//   - `"hoeffding"`, `"hoeffd"`, `"d"`: Hoeffding's D
func Wdm(x, y []float64, method string) float64 {
	if len(x) == 0 || len(y) == 0 {
		return 0
	}

	return wrapper.Wdm(&x[0], len(x), &y[0], len(y), method)
}
