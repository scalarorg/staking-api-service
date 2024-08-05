package utils

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/babylonchain/babylon/btcstaking"
	"github.com/babylonchain/babylon/crypto/bip322"
	bbntypes "github.com/babylonchain/babylon/types"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"

	"github.com/scalarorg/staking-api-service/internal/types"
)

// GetSchnorrPkFromHex parses Schnorr public keys in 32 bytes
func GetSchnorrPkFromHex(pkHex string) (*btcec.PublicKey, error) {
	pkBytes, err := hex.DecodeString(pkHex)
	if err != nil {
		return nil, err
	}

	return schnorr.ParsePubKey(pkBytes)
}

// GetCovenantPksFromStrings parses BTC public keys in 33 bytes
func GetCovenantPksFromStrings(pkStrings []string) ([]*btcec.PublicKey, error) {
	pks := make([]*btcec.PublicKey, len(pkStrings))
	for i, pkStr := range pkStrings {
		pkBytes, err := hex.DecodeString(pkStr)
		if err != nil {
			return nil, err
		}

		pk, err := btcec.ParsePubKey(pkBytes)
		if err != nil {
			return nil, err
		}

		pks[i] = pk
	}

	return pks, nil
}

func parseUnbondingTxHex(unbondingTxHex string) (*wire.MsgTx, error) {
	unbondingTx, _, err := bbntypes.NewBTCTxFromHex(unbondingTxHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode unbonding tx from hex: %w", err)
	}

	if err := btcstaking.IsSimpleTransfer(unbondingTx); err != nil {
		return nil, fmt.Errorf("the unbonding tx is not a simple transfer: %w", err)
	}

	return unbondingTx, nil
}

func VerifyUnbondingRequest(
	stakingTxHashHex,
	unbondingTxHashHex,
	unbondingTxHex,
	stakerPkHex,
	finalityProviderPkHex,
	unbondingSigHex string,
	stakingTimeLock,
	stakingOutputIndex,
	stakingValue uint64,
	params *types.VersionedGlobalParams,
	btcNetParam *chaincfg.Params,
) error {
	// 1. validate that un-bonding transaction has proper shape
	unbondingTx, err := parseUnbondingTxHex(unbondingTxHex)
	if err != nil {
		return fmt.Errorf("invalid unbonding tx hex: %w", err)
	}

	// 2. validate that un-bonding tx hash is valid and matches the hash of the
	// provided unbonding tx
	unbondingTxHash, err := chainhash.NewHashFromStr(unbondingTxHashHex)

	if err != nil {
		return fmt.Errorf("failed to decode unbonding tx hash from hex: %w", err)
	}

	unbondingTxHashFromTx := unbondingTx.TxHash()

	if !unbondingTxHashFromTx.IsEqual(unbondingTxHash) {
		return fmt.Errorf("unbonding_tx_hash_hex must match the hash calculated from the provided unbonding tx")
	}

	// 3. validate the un-bonding transaction points to the previous staking tx
	stakingTxHash, err := chainhash.NewHashFromStr(stakingTxHashHex)
	if err != nil {
		return fmt.Errorf("failed to decode staking tx hash from hex: %w", err)
	}
	if !unbondingTx.TxIn[0].PreviousOutPoint.Hash.IsEqual(stakingTxHash) {
		return fmt.Errorf("the unbonding tx input must match the previous staking tx hash, expected: %s, got: %s",
			stakingTxHashHex,
			unbondingTx.TxIn[0].PreviousOutPoint.Hash.String(),
		)
	}
	if uint64(unbondingTx.TxIn[0].PreviousOutPoint.Index) != stakingOutputIndex {
		return fmt.Errorf("the unbonding tx input must match the previous staking tx output index, expected: %d, got: %d",
			stakingOutputIndex,
			unbondingTx.TxIn[0].PreviousOutPoint.Index,
		)
	}

	// 4. verify that the unbonding output is constructed as expected
	covenantPks, err := GetCovenantPksFromStrings(params.CovenantPks)
	if err != nil {
		return fmt.Errorf("failed to decode coveant public keys from strings: %w", err)
	}

	stakerPk, err := GetSchnorrPkFromHex(stakerPkHex)
	if err != nil {
		return fmt.Errorf("failed to decode staker public key from hex: %w", err)
	}

	finalityProviderPk, err := GetSchnorrPkFromHex(finalityProviderPkHex)
	if err != nil {
		return fmt.Errorf("failed to decode finality provider public key from hex: %w", err)
	}

	expectedUnbondingOutputValue := btcutil.Amount(stakingValue) - btcutil.Amount(params.UnbondingFee)
	if expectedUnbondingOutputValue <= 0 {
		return fmt.Errorf("staking output value is too low, got %v, unbonding fee: %v",
			btcutil.Amount(stakingValue), btcutil.Amount(params.UnbondingFee))
	}

	unbondingInfo, err := btcstaking.BuildUnbondingInfo(
		stakerPk,
		[]*btcec.PublicKey{finalityProviderPk},
		covenantPks,
		uint32(params.CovenantQuorum),
		uint16(params.UnbondingTime),
		expectedUnbondingOutputValue,
		btcNetParam,
	)
	if err != nil {
		return fmt.Errorf("failed to build unbonding info")
	}

	if !outputsAreEqual(unbondingInfo.UnbondingOutput, unbondingTx.TxOut[0]) {
		return fmt.Errorf("unbonding output does not match expected output")
	}

	// 5. verify the signature
	stakingInfo, err := btcstaking.BuildStakingInfo(
		stakerPk,
		[]*btcec.PublicKey{finalityProviderPk},
		covenantPks,
		uint32(params.CovenantQuorum),
		uint16(stakingTimeLock),
		btcutil.Amount(stakingValue),
		btcNetParam,
	)
	if err != nil {
		return fmt.Errorf("failed to build staking info")
	}
	sigBytes, err := hex.DecodeString(unbondingSigHex)
	if err != nil {
		return fmt.Errorf("failed to decode unbonding signature from hex")
	}
	unbondingSpendInfo, err := stakingInfo.UnbondingPathSpendInfo()
	if err != nil {
		return fmt.Errorf("failed to build unbonding path spend info")
	}
	if err := btcstaking.VerifyTransactionSigWithOutput(
		unbondingTx,
		stakingInfo.StakingOutput,
		unbondingSpendInfo.GetPkScriptPath(),
		stakerPk,
		sigBytes,
	); err != nil {
		return fmt.Errorf("invalid unbonding signature")
	}
	return nil
}

func outputsAreEqual(a *wire.TxOut, b *wire.TxOut) bool {
	if a.Value != b.Value {
		return false
	}

	if !bytes.Equal(a.PkScript, b.PkScript) {
		return false
	}

	return true
}

func GetBtcNetParamesFromString(net string) (*chaincfg.Params, error) {
	var netParams chaincfg.Params
	switch net {
	case "mainnet":
		netParams = chaincfg.MainNetParams
	case "testnet3":
		netParams = chaincfg.TestNet3Params
	case "regtest":
		netParams = chaincfg.RegressionNetParams
	case "simnet":
		netParams = chaincfg.SimNetParams
	case "signet":
		netParams = chaincfg.SigNetParams
	default:
		return nil, fmt.Errorf("invalid network: %s", net)
	}
	return &netParams, nil
}

func GetTaprootAddressFromPk(pkHex string, netParams *chaincfg.Params) (string, error) {
	pk, err := GetSchnorrPkFromHex(pkHex)
	if err != nil {
		return "", err
	}
	address, err := bip322.PubKeyToP2TrSpendAddress(pk, netParams)
	if err != nil {
		return "", err
	}
	return address.EncodeAddress(), nil
}
