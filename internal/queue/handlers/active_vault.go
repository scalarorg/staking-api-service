package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/staking-api-service/internal/types"
	queueClient "github.com/scalarorg/staking-queue-client/client"
)

// ActiveVaultHandler handles the active vault event
// This handler is designed to be idempotent, capable of handling duplicate messages gracefully.
// It can also resume from the next step if a previous step fails, ensuring robustness in the event processing workflow.
func (h *QueueHandler) ActiveVaultHandler(ctx context.Context, messageBody string) *types.Error {
	// Parse the message body into ActiveVaultEvent
	var activeVaultEvent queueClient.ActiveVaultEvent
	err := json.Unmarshal([]byte(messageBody), &activeVaultEvent)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to unmarshal the message body into ActiveVaultEvent")
		return types.NewError(http.StatusBadRequest, types.BadRequest, err)
	}

	// Check if delegation already exists
	exist, delError := h.Services.IsDelegationPresent(ctx, activeVaultEvent.VaultTxHashHex)
	if delError != nil {
		return delError
	}
	if exist {
		// Ignore the message as the delegation already exists. This is a duplicate message
		log.Ctx(ctx).Debug().Str("VaultTxHashHex", activeVaultEvent.VaultTxHashHex).
			Msg("delegation already exists")
		return nil
	}

	// We only emit the stats event for the active vault if it is not an overflow event
	if !activeVaultEvent.IsOverflow {
		// Perform the async stats calculation by emit the stats event
		statsError := h.EmitStatsEvent(ctx, queueClient.NewStatsEvent(
			activeVaultEvent.VaultTxHashHex,
			activeVaultEvent.StakerPkHex,
			activeVaultEvent.FinalityProviderPkHex,
			activeVaultEvent.StakingValue,
			types.Active.ToString(),
		))
		if statsError != nil {
			log.Ctx(ctx).Error().Err(statsError).Msg("Failed to emit stats event for active vault")
			return statsError
		}
	}

	// Save the active vault delegation. This is the final step in the active vault event processing
	// Please refer to the README.md for the details on the active vault event processing workflow
	// This function not dependent on any other function and can be executed independently
	// so we can save active vault delegation cause it have all the required field
	saveErr := h.Services.SaveActiveStakingDelegation(
		ctx, activeVaultEvent.VaultTxHashHex, activeVaultEvent.StakerPkHex,
		activeVaultEvent.FinalityProviderPkHex, activeVaultEvent.StakingValue,
		activeVaultEvent.StakingStartHeight, activeVaultEvent.StakingStartTimestamp,
		0, activeVaultEvent.StakingOutputIndex,
		activeVaultEvent.VaultTxHex, activeVaultEvent.IsOverflow,
	)
	if saveErr != nil {
		return saveErr
	}

	return nil
}
