package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/staking-api-service/internal/types"
	"github.com/scalarorg/staking-api-service/internal/utils"
	queueClient "github.com/scalarorg/staking-queue-client/client"
)

func (h *QueueHandler) BurningVaultHandler(ctx context.Context, messageBody string) *types.Error {
	var burningVaultEvent queueClient.BurningVaultEvent
	err := json.Unmarshal([]byte(messageBody), &burningVaultEvent)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to unmarshal the message body into burningVaultEvent")
		return types.NewError(http.StatusBadRequest, types.BadRequest, err)
	}

	// Check if the delegation is in the right state to process the burning event.
	del, delErr := h.Services.GetDelegation(ctx, burningVaultEvent.VaultTxHashHex)
	// Requeue if found any error. Including not found error
	if delErr != nil {
		return delErr
	}
	state := del.State

	VaultTxHashHex := burningVaultEvent.GetVaultTxHashHex()

	if utils.Contains(utils.OutdatedStatesForBurning(), state) {
		// Ignore the message as the delegation state is burning.Nothing to do anymore
		log.Ctx(ctx).Debug().Str("vaultTxHashHex", VaultTxHashHex).
			Msg("delegation state is outdated for burning event")
		return nil
	}
	// Requeue if the current state is not in the qualified states to transition to burning
	// We will wait for the unbonded message to be processed first.
	if !utils.Contains(utils.QualifiedStatesToBurning(), state) {
		errMsg := "delegation is not in the qualified state to transition to burning"
		log.Ctx(ctx).Warn().Str("vaultTxHashHex", VaultTxHashHex).
			Str("state", state.ToString()).Msg(errMsg)
		return types.NewErrorWithMsg(http.StatusForbidden, types.Forbidden, errMsg)
	}

	// Transition to burning state
	// Please refer to the README.md for the details on the event processing workflow
	transitionErr := h.Services.TransitionToBurningState(
		ctx, burningVaultEvent.VaultTxHashHex,
	)
	if transitionErr != nil {
		return transitionErr
	}

	return nil
}

func (h *QueueHandler) SlashingOrLostKeyVaultHandler(ctx context.Context, messageBody string) *types.Error {
	var SlashingOrLostKeyVaultEvent queueClient.SlashingOrLostKeyVaultEvent
	err := json.Unmarshal([]byte(messageBody), &SlashingOrLostKeyVaultEvent)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to unmarshal the message body into SlashingOrLostKeyVaultEvent")
		return types.NewError(http.StatusBadRequest, types.BadRequest, err)
	}

	// Check if the delegation is in the right state to process the burning event.
	del, delErr := h.Services.GetDelegation(ctx, SlashingOrLostKeyVaultEvent.VaultTxHashHex)
	// Requeue if found any error. Including not found error
	if delErr != nil {
		return delErr
	}
	state := del.State

	VaultTxHashHex := SlashingOrLostKeyVaultEvent.GetVaultTxHashHex()

	if utils.Contains(utils.OutdatedStatesForSlashingOrLostKey(), state) {
		// Ignore the message as the delegation state is burning.         uy88Nothing to do anymore
		log.Ctx(ctx).Debug().Str("vaultTxHashHex", VaultTxHashHex).
			Msg("delegation state is outdated for burning event")
		return nil
	}
	// Requeue if the current state is not in the qualified states to transition to burning
	// We will wait for the unbonded message to be processed first.
	if !utils.Contains(utils.QualifiedStatesToSlashingOrLostKey(), state) {
		errMsg := "delegation is not in the qualified state to transition to burning"
		log.Ctx(ctx).Warn().Str("vaultTxHashHex", VaultTxHashHex).
			Str("state", state.ToString()).Msg(errMsg)
		return types.NewErrorWithMsg(http.StatusForbidden, types.Forbidden, errMsg)
	}

	// Transition to burning state
	// Please refer to the README.md for the details on the event processing workflow
	transitionErr := h.Services.TransitionToSlashingOrLostKeyState(
		ctx, SlashingOrLostKeyVaultEvent.VaultTxHashHex,
	)
	if transitionErr != nil {
		return transitionErr
	}

	return nil
}

func (h *QueueHandler) BurnWithoutDAppVaultHandler(ctx context.Context, messageBody string) *types.Error {
	var BurnWithoutDAppVaultEvent queueClient.BurnWithoutDAppVaultEvent
	err := json.Unmarshal([]byte(messageBody), &BurnWithoutDAppVaultEvent)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to unmarshal the message body into BurnWithoutDAppVaultEvent")
		return types.NewError(http.StatusBadRequest, types.BadRequest, err)
	}

	// Check if the delegation is in the right state to process the burning event.
	del, delErr := h.Services.GetDelegation(ctx, BurnWithoutDAppVaultEvent.VaultTxHashHex)
	// Requeue if found any error. Including not found error
	if delErr != nil {
		return delErr
	}
	state := del.State

	VaultTxHashHex := BurnWithoutDAppVaultEvent.GetVaultTxHashHex()

	if utils.Contains(utils.OutdatedStatesForBurnWithoutDApp(), state) {
		// Ignore the message as the delegation state is burning.         uy88Nothing to do anymore
		log.Ctx(ctx).Debug().Str("vaultTxHashHex", VaultTxHashHex).
			Msg("delegation state is outdated for burning event")
		return nil
	}
	// Requeue if the current state is not in the qualified states to transition to burning
	// We will wait for the unbonded message to be processed first.
	if !utils.Contains(utils.QualifiedStatesToBurnWithoutDApp(), state) {
		errMsg := "delegation is not in the qualified state to transition to burning"
		log.Ctx(ctx).Warn().Str("vaultTxHashHex", VaultTxHashHex).
			Str("state", state.ToString()).Msg(errMsg)
		return types.NewErrorWithMsg(http.StatusForbidden, types.Forbidden, errMsg)
	}

	// Transition to burning state
	// Please refer to the README.md for the details on the event processing workflow
	transitionErr := h.Services.TransitionToBurnWithoutDAppState(
		ctx, BurnWithoutDAppVaultEvent.VaultTxHashHex,
	)
	if transitionErr != nil {
		return transitionErr
	}

	return nil
}
