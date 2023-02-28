package services_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yvv4git/task-wow/internal/services"
	"github.com/yvv4git/task-wow/pkg/pow"
)

func TestWOW(t *testing.T) {
	ctx := context.Background()

	// Step-1: after client connected server send challenge to client.
	sha256PoW := pow.NewSHA256(services.DefaultDifficulty)
	wowSvc := services.NewWOW(sha256PoW)
	challengeMsg, err := wowSvc.Send(ctx)
	require.NoError(t, err)
	t.Logf("challengeMsg to client: %v", challengeMsg)

	// Step-2: client receives the challenge and sends the solution to the server.
	proof := sha256PoW.SolveChallenge([]byte(challengeMsg))
	t.Logf("proff: %#v", proof)

	rawMsg, err := json.Marshal(proof)
	require.NoError(t, err)
	t.Logf("raw msg: %v", string(rawMsg))

	// Step-3: server receives the challenge solution and check this.
	var proofFromClient pow.Proof
	err = json.Unmarshal(rawMsg, &proofFromClient)
	require.NoError(t, err)

	statusCheck := sha256PoW.Check([]byte(challengeMsg), proofFromClient)
	t.Logf("status check: %v", statusCheck)
}
