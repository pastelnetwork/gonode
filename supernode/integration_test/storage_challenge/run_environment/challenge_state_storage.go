package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type challengeStateStorage struct{}

func saveState(state string, challengeID, nodeID string, sentBlock int32) {
	formval := url.Values{}
	formval.Add("id", challengeID)
	formval.Add("node_id", nodeID)
	formval.Add("sent_block", fmt.Sprint(sentBlock))
	http.DefaultClient.PostForm(fmt.Sprintf("http://helper:8088/sts/%s", state), formval)
}

func (*challengeStateStorage) OnSent(_ context.Context, challengeID, nodeID string, sentBlock int32) {
	saveState("sent", challengeID, nodeID, sentBlock)
}

func (*challengeStateStorage) OnResponded(_ context.Context, challengeID, nodeID string, sentBlock int32) {
	saveState("respond", challengeID, nodeID, sentBlock)
}

func (*challengeStateStorage) OnSucceeded(_ context.Context, challengeID, nodeID string, sentBlock int32) {
	saveState("succeeded", challengeID, nodeID, sentBlock)
}

func (*challengeStateStorage) OnFailed(_ context.Context, challengeID, nodeID string, sentBlock int32) {
	saveState("failed", challengeID, nodeID, sentBlock)
}

func (*challengeStateStorage) OnTimeout(_ context.Context, challengeID, nodeID string, sentBlock int32) {
	saveState("timeout", challengeID, nodeID, sentBlock)
}
