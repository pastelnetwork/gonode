package storagechallenge

// this is using for testing

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"net/url"

// 	"github.com/pastelnetwork/gonode/common/log"
// )

// var httpClient = http.DefaultClient

// func (s *service) sendStatictis(ctx context.Context, challengeMsg *ChallengeMessage, path string) {
// 	formVal := url.Values{}
// 	formVal.Add("id", challengeMsg.ChallengeID)
// 	formVal.Add("node_id", challengeMsg.ChallengingMasternodeID)
// 	formVal.Add("sent_block", fmt.Sprint(challengeMsg.BlockNumChallengeSent))
// 	formVal.Add("key", challengeMsg.FileHashToChallenge)
// 	httpClient.PostForm(fmt.Sprintf("http://helper:8088/sts/%s", path), formVal)
// 	log.WithContext(ctx).WithField("challengeID", challengeMsg.ChallengeID).Infof("query to update %s statictis data", path)
// }
