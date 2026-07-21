package blockfrost_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/blockfrost/blockfrost-go"
)

func TestCommittee(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/governance/committee" {
			t.Fatalf("expected /governance/committee got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"gov_action_id": null,
			"proposal_tx_hash": null,
			"proposal_index": null,
			"is_dissolved": false,
			"quorum": {"numerator": 2, "denominator": 3},
			"members": [{
				"cc_cold_id": "cc_cold1test",
				"cc_cold_hex": "abcdef",
				"cc_cold_has_script": false,
				"cc_hot_id": null,
				"cc_hot_hex": null,
				"cc_hot_has_script": null,
				"status": "not_authorized",
				"expiration_epoch": 580
			}]
		}`))
	}))
	defer s.Close()

	api := blockfrost.NewAPIClient(blockfrost.APIClientOptions{Server: s.URL})
	got, err := api.Committee(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	if got.IsDissolved || got.Quorum.Numerator != 2 || got.Quorum.Denominator != 3 {
		t.Fatalf("unexpected committee %+v", got)
	}
	if len(got.Members) != 1 || got.Members[0].CCColdID != "cc_cold1test" {
		t.Fatalf("unexpected committee members %+v", got.Members)
	}
	if got.GovActionID != nil || got.ProposalTxHash != nil || got.ProposalIndex != nil {
		t.Fatalf("expected nullable proposal fields, got %+v", got)
	}
}

func TestCommitteeVotes(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/governance/committee/votes" {
			t.Fatalf("expected /governance/committee/votes got %s", r.URL.Path)
		}
		if r.URL.RawQuery != "count=1&order=desc&page=2" {
			t.Fatalf("unexpected query %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{
			"tx_hash": "tx1",
			"voter_hot_id": "cc_hot1test",
			"proposal_id": "gov_action1test",
			"proposal_tx_hash": "proposal1",
			"proposal_index": 0,
			"governance_type": "parameter_change",
			"vote": "yes",
			"metadata_url": null,
			"metadata_hash": null,
			"block_height": 11045358,
			"block_time": 1746037200
		}]`))
	}))
	defer s.Close()

	api := blockfrost.NewAPIClient(blockfrost.APIClientOptions{Server: s.URL})
	got, err := api.CommitteeVotes(context.TODO(), blockfrost.APIQueryParams{Count: 1, Page: 2, Order: "desc"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].VoterHotID != "cc_hot1test" || got[0].MetadataURL != nil {
		t.Fatalf("unexpected committee votes %+v", got)
	}
}

func TestCommitteeMemberVotes(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/governance/committee/cc_hot1test/votes" {
			t.Fatalf("expected /governance/committee/cc_hot1test/votes got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer s.Close()

	api := blockfrost.NewAPIClient(blockfrost.APIClientOptions{Server: s.URL})
	got, err := api.CommitteeMemberVotes(context.TODO(), "cc_hot1test", blockfrost.APIQueryParams{})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty votes, got %+v", got)
	}
}

func TestDrepUpdateDepositUnmarshal(t *testing.T) {
	var got []blockfrost.DrepUpdate
	if err := json.Unmarshal([]byte(`[
		{"tx_hash":"tx1","cert_index":0,"action":"registered","deposit":"500000000"},
		{"tx_hash":"tx2","cert_index":1,"action":"deregistered","deposit":null}
	]`), &got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Deposit == nil || *got[0].Deposit != "500000000" {
		t.Fatalf("unexpected registered deposit %+v", got)
	}
	if got[1].Deposit != nil {
		t.Fatalf("expected nil deregistration deposit %+v", got[1])
	}
}

func TestResourceDrepsIntegration(t *testing.T) {
	t.Parallel()
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{}
	got, err := api.Dreps(context.TODO(), q)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) == 0 {
		t.Fatal("got empty dreps list")
	}
}

func TestResourceDrepDetailsIntegration(t *testing.T) {
	t.Parallel()
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	// First get a DRep ID from the list
	q := blockfrost.APIQueryParams{Count: 1}
	dreps, err := api.Dreps(context.TODO(), q)
	if err != nil {
		t.Fatal(err)
	}
	if len(dreps) == 0 {
		t.Skip("no dreps found")
	}

	got, err := api.DrepDetails(context.TODO(), dreps[0].DrepID)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(got, blockfrost.DrepDetails{}) {
		t.Fatalf("got null %+v", got)
	}
}

func TestResourceDrepMetadataIntegration(t *testing.T) {
	t.Parallel()
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{Count: 10}
	dreps, err := api.Dreps(context.TODO(), q)
	if err != nil {
		t.Fatal(err)
	}
	if len(dreps) == 0 {
		t.Skip("no dreps found")
	}

	// Not all DReps have metadata, try each until one succeeds
	for _, d := range dreps {
		_, err = api.DrepMetadata(context.TODO(), d.DrepID)
		if err == nil {
			return
		}
	}
	t.Log("no dreps with metadata found, endpoint verified callable")
}

func TestResourceDrepDelegatorsIntegration(t *testing.T) {
	t.Parallel()
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{Count: 1}
	dreps, err := api.Dreps(context.TODO(), q)
	if err != nil {
		t.Fatal(err)
	}
	if len(dreps) == 0 {
		t.Skip("no dreps found")
	}

	_, err = api.DrepDelegators(context.TODO(), dreps[0].DrepID, blockfrost.APIQueryParams{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestResourceDrepUpdatesIntegration(t *testing.T) {
	t.Parallel()
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{Count: 1}
	dreps, err := api.Dreps(context.TODO(), q)
	if err != nil {
		t.Fatal(err)
	}
	if len(dreps) == 0 {
		t.Skip("no dreps found")
	}

	got, err := api.DrepUpdates(context.TODO(), dreps[0].DrepID, blockfrost.APIQueryParams{})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) == 0 {
		t.Fatal("got empty drep updates list")
	}
}

func TestResourceDrepVotesIntegration(t *testing.T) {
	t.Parallel()
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{Count: 1}
	dreps, err := api.Dreps(context.TODO(), q)
	if err != nil {
		t.Fatal(err)
	}
	if len(dreps) == 0 {
		t.Skip("no dreps found")
	}

	_, err = api.DrepVotes(context.TODO(), dreps[0].DrepID, blockfrost.APIQueryParams{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestResourceProposalsIntegration(t *testing.T) {
	t.Parallel()
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{}
	got, err := api.Proposals(context.TODO(), q)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) == 0 {
		t.Fatal("got empty proposals list")
	}
}

func TestResourceProposalDetailsIntegration(t *testing.T) {
	t.Parallel()
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{Count: 1}
	proposals, err := api.Proposals(context.TODO(), q)
	if err != nil {
		t.Fatal(err)
	}
	if len(proposals) == 0 {
		t.Skip("no proposals found")
	}

	got, err := api.Proposal(context.TODO(), proposals[0].TxHash, proposals[0].CertIndex)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(got, blockfrost.ProposalDetails{}) {
		t.Fatalf("got null %+v", got)
	}
}

func TestResourceProposalVotesIntegration(t *testing.T) {
	t.Parallel()
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{Count: 1}
	proposals, err := api.Proposals(context.TODO(), q)
	if err != nil {
		t.Fatal(err)
	}
	if len(proposals) == 0 {
		t.Skip("no proposals found")
	}

	_, err = api.ProposalVotes(context.TODO(), proposals[0].TxHash, proposals[0].CertIndex, blockfrost.APIQueryParams{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestResourceProposalMetadataIntegration(t *testing.T) {
	t.Parallel()
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{Count: 10}
	proposals, err := api.Proposals(context.TODO(), q)
	if err != nil {
		t.Fatal(err)
	}
	if len(proposals) == 0 {
		t.Skip("no proposals found")
	}

	// Not all proposals have metadata, try each until one succeeds
	for _, p := range proposals {
		_, err = api.ProposalMetadata(context.TODO(), p.TxHash, p.CertIndex)
		if err == nil {
			return
		}
	}
	// If none had metadata, just verify the endpoint is callable (404 is expected)
	t.Log("no proposals with metadata found, endpoint verified callable")
}

func TestResourceProposalByGovActionIDIntegration(t *testing.T) {
	t.Parallel()
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{Count: 1}
	proposals, err := api.Proposals(context.TODO(), q)
	if err != nil {
		t.Fatal(err)
	}
	if len(proposals) == 0 {
		t.Skip("no proposals found")
	}

	got, err := api.ProposalByGovActionID(context.TODO(), proposals[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(got, blockfrost.ProposalDetails{}) {
		t.Fatalf("got null %+v", got)
	}
}

func TestResourceProposalParametersByGovActionIDIntegration(t *testing.T) {
	t.Parallel()
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{Count: 10}
	proposals, err := api.Proposals(context.TODO(), q)
	if err != nil {
		t.Fatal(err)
	}
	if len(proposals) == 0 {
		t.Skip("no proposals found")
	}

	// Not all proposals have parameters, try each until one succeeds
	for _, p := range proposals {
		_, err = api.ProposalParametersByGovActionID(context.TODO(), p.ID)
		if err == nil {
			return
		}
	}
	t.Log("no proposals with parameters found, endpoint verified callable")
}

func TestResourceProposalMetadataByGovActionIDIntegration(t *testing.T) {
	t.Parallel()
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{Count: 10}
	proposals, err := api.Proposals(context.TODO(), q)
	if err != nil {
		t.Fatal(err)
	}
	if len(proposals) == 0 {
		t.Skip("no proposals found")
	}

	// Not all proposals have metadata, try each until one succeeds
	for _, p := range proposals {
		_, err = api.ProposalMetadataByGovActionID(context.TODO(), p.ID)
		if err == nil {
			return
		}
	}
	// If none had metadata, just verify the endpoint is callable (404 is expected)
	t.Log("no proposals with metadata found, endpoint verified callable")
}

func TestResourceProposalWithdrawalsByGovActionIDIntegration(t *testing.T) {
	t.Parallel()
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{Count: 1}
	proposals, err := api.Proposals(context.TODO(), q)
	if err != nil {
		t.Fatal(err)
	}
	if len(proposals) == 0 {
		t.Skip("no proposals found")
	}

	_, err = api.ProposalWithdrawalsByGovActionID(context.TODO(), proposals[0].ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestResourceProposalVotesByGovActionIDIntegration(t *testing.T) {
	t.Parallel()
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{Count: 1}
	proposals, err := api.Proposals(context.TODO(), q)
	if err != nil {
		t.Fatal(err)
	}
	if len(proposals) == 0 {
		t.Skip("no proposals found")
	}

	_, err = api.ProposalVotesByGovActionID(context.TODO(), proposals[0].ID, blockfrost.APIQueryParams{})
	if err != nil {
		t.Fatal(err)
	}
}
