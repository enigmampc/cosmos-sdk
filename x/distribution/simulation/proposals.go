package simulation

import (
	"math/rand"

	simappparams "github.com/enigmampc/cosmos-sdk/simapp/params"
	sdk "github.com/enigmampc/cosmos-sdk/types"
	"github.com/enigmampc/cosmos-sdk/x/distribution/keeper"
	"github.com/enigmampc/cosmos-sdk/x/distribution/types"
	govtypes "github.com/enigmampc/cosmos-sdk/x/gov/types"
	"github.com/enigmampc/cosmos-sdk/x/simulation"
)

// Proposal operation weights
const (
	OpWeightSubmitCommunitySpendProposal      = "op_weight_submit_community_spend_proposal"
	OpWeightSubmitSecretFoundationTaxProposal = "op_weight_submit_secret_foundation_tax_proposal"
)

// ProposalContents defines the module weighted proposals' contents
func ProposalContents(k keeper.Keeper) []simulation.WeightedProposalContent {
	return []simulation.WeightedProposalContent{
		{
			AppParamsKey:       OpWeightSubmitCommunitySpendProposal,
			DefaultWeight:      simappparams.DefaultWeightCommunitySpendProposal,
			ContentSimulatorFn: SimulateCommunityPoolSpendProposalContent(k),
		},
		{
			AppParamsKey:       OpWeightSubmitSecretFoundationTaxProposal,
			DefaultWeight:      simappparams.DefaultWeightSecretFoundationTaxProposal,
			ContentSimulatorFn: SimulateSecretFoundationTaxProposalContent(),
		},
	}
}

// SimulateCommunityPoolSpendProposalContent generates random community-pool-spend proposal content
// nolint: funlen
func SimulateCommunityPoolSpendProposalContent(k keeper.Keeper) simulation.ContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, accs []simulation.Account) govtypes.Content {
		simAccount, _ := simulation.RandomAcc(r, accs)

		balance := k.GetFeePool(ctx).CommunityPool
		if balance.Empty() {
			return nil
		}

		denomIndex := r.Intn(len(balance))
		amount, err := simulation.RandPositiveInt(r, balance[denomIndex].Amount.TruncateInt())
		if err != nil {
			return nil
		}

		return types.NewCommunityPoolSpendProposal(
			simulation.RandStringOfLength(r, 10),
			simulation.RandStringOfLength(r, 100),
			simAccount.Address,
			sdk.NewCoins(sdk.NewCoin(balance[denomIndex].Denom, amount)),
		)
	}
}

// SimulateSecretFoundationTaxProposalContent generates a random SecretFoundationTaxProposal.
func SimulateSecretFoundationTaxProposalContent() simulation.ContentSimulatorFn {
	return func(r *rand.Rand, ctx sdk.Context, accs []simulation.Account) govtypes.Content {
		simAccount, _ := simulation.RandomAcc(r, accs)

		return types.NewSecretFoundationTaxProposal(
			simulation.RandStringOfLength(r, 10),
			simulation.RandStringOfLength(r, 100),
			simulation.RandomDecAmount(r, sdk.MustNewDecFromStr("0.20")),
			simAccount.Address,
		)
	}
}
