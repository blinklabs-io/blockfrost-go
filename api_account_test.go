package blockfrost_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/blockfrost/blockfrost-go"
)

func TestAccountRegistrationHistoryDepositUnmarshal(t *testing.T) {
	var got []blockfrost.AccountRegistrationHistory
	if err := json.Unmarshal([]byte(`[
		{"tx_hash":"tx1","action":"registered","deposit":"2000000","tx_slot":45093580,"block_time":1646437200,"block_height":6745358},
		{"tx_hash":"tx2","action":"deregistered","deposit":null,"tx_slot":48093580,"block_time":1649033600,"block_height":7126896}
	]`), &got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Deposit == nil || *got[0].Deposit != "2000000" {
		t.Fatalf("unexpected registered deposit %+v", got)
	}
	if got[1].Deposit != nil {
		t.Fatalf("expected nil deregistration deposit %+v", got[1])
	}
}

func TestResourceAccountIntegration(t *testing.T) {
	t.Parallel()

	inputStakeAddr := "stake1ux3g2c9dx2nhhehyrezyxpkstartcqmu9hk63qgfkccw5rqttygt7"
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	got, err := api.Account(context.TODO(), inputStakeAddr)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if reflect.DeepEqual(got, blockfrost.Account{}) {
		t.Fatalf("got null %+v", got)
	}
}

func TestResourceAccountRewardsHistoryIntegration(t *testing.T) {
	t.Parallel()
	inputStakeAddr := "stake1ux3g2c9dx2nhhehyrezyxpkstartcqmu9hk63qgfkccw5rqttygt7"
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{
		Count: 1,
	}
	got, err := api.AccountRewardsHistory(context.TODO(), inputStakeAddr, q)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if reflect.DeepEqual(got, []blockfrost.AccountRewardsHistory{}) {
		t.Fatalf("got null %+v", got)
	}

}

func TestResourceAccountHistoryIntegration(t *testing.T) {
	t.Parallel()

	inputStakeAddr := "stake1ux3g2c9dx2nhhehyrezyxpkstartcqmu9hk63qgfkccw5rqttygt7"
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{
		Count: 1,
	}
	got, err := api.AccountHistory(context.TODO(), inputStakeAddr, q)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if reflect.DeepEqual(got, []blockfrost.AccountHistory{}) {
		t.Fatalf("got null %+v", got)
	}

}

func TestResourceAccountDelegationHistoryIntregration(t *testing.T) {
	t.Parallel()

	inputStakeAddr := "stake1ux3g2c9dx2nhhehyrezyxpkstartcqmu9hk63qgfkccw5rqttygt7"
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)
	q := blockfrost.APIQueryParams{
		Count: 1,
	}
	got, err := api.AccountDelegationHistory(context.TODO(), inputStakeAddr, q)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if reflect.DeepEqual(got, []blockfrost.AccountDelegationHistory{}) {
		t.Fatalf("got null %+v", got)
	}

}

func TestResourceAccountRegistrationHistoryIntegration(t *testing.T) {
	t.Parallel()
	inputStakeAddr := "stake1ux3g2c9dx2nhhehyrezyxpkstartcqmu9hk63qgfkccw5rqttygt7"
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{
		Count: 1,
	}
	got, err := api.AccountRegistrationHistory(context.TODO(), inputStakeAddr, q)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if reflect.DeepEqual(got, []blockfrost.AccountRegistrationHistory{}) {
		t.Fatalf("got null %+v", got)
	}
}

func TestResourceAccountWithdrawalHistoryIntegration(t *testing.T) {
	t.Parallel()
	inputStakeAddr := "stake1ux3g2c9dx2nhhehyrezyxpkstartcqmu9hk63qgfkccw5rqttygt7"
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{
		Count: 1,
	}
	got, err := api.AccountWithdrawalHistory(context.TODO(), inputStakeAddr, q)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if reflect.DeepEqual(got, []blockfrost.AccountWithdrawalHistory{}) {
		t.Fatalf("got null %+v", got)
	}
}

func TestResourceAccountMIRHistoryIntegration(t *testing.T) {
	t.Parallel()
	inputStakeAddr := "stake1ux3g2c9dx2nhhehyrezyxpkstartcqmu9hk63qgfkccw5rqttygt7"
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{
		Count: 1,
	}
	got, err := api.AccountMIRHistory(context.TODO(), inputStakeAddr, q)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if reflect.DeepEqual(got, []blockfrost.AccountMIRHistory{}) {
		t.Fatalf("got null %+v", got)
	}
}

func TestResourceAccountAssociatedAddressesIntegration(t *testing.T) {
	t.Parallel()
	inputStakeAddr := "stake1ux3g2c9dx2nhhehyrezyxpkstartcqmu9hk63qgfkccw5rqttygt7"
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{
		Count: 1,
	}
	got, err := api.AccountAssociatedAddresses(context.TODO(), inputStakeAddr, q)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if reflect.DeepEqual(got, []blockfrost.AccountAssociatedAddress{}) {
		t.Fatalf("got null %+v", got)
	}
}

func TestResourceAccountAssociatedAssetsIntegration(t *testing.T) {
	t.Parallel()

	inputStakeAddr := "stake1u9ylzsgxaa6xctf4juup682ar3juj85n8tx3hthnljg47zctvm3rc"
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{
		Count: 1,
	}
	_, err := api.AccountAssociatedAssets(context.TODO(), inputStakeAddr, q)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestResourceAccountAddressTotalIntegration(t *testing.T) {
	t.Parallel()
	inputStakeAddr := "stake1ux3g2c9dx2nhhehyrezyxpkstartcqmu9hk63qgfkccw5rqttygt7"
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	got, err := api.AccountAddressesTotal(context.TODO(), inputStakeAddr)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(got, blockfrost.AccountAddressesTotal{}) {
		t.Fatalf("got null %+v", got)
	}
}

func TestResourceAccountTransactionsIntegration(t *testing.T) {
	t.Parallel()
	inputStakeAddr := "stake1ux3g2c9dx2nhhehyrezyxpkstartcqmu9hk63qgfkccw5rqttygt7"
	api := blockfrost.NewAPIClient(
		blockfrost.APIClientOptions{},
	)

	q := blockfrost.APIQueryParams{Count: 1}
	got, err := api.AccountTransactions(context.TODO(), inputStakeAddr, q)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) == 0 {
		t.Fatal("got empty account transactions")
	}
}
