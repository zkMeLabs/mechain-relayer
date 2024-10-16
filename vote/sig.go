package vote

import (
	"encoding/hex"
	"reflect"

	"github.com/0xPolygon/polygon-edge/bls"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cometbft/cometbft/votepool"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/willf/bitset"

	"github.com/zkMeLabs/mechain-relayer/db/model"
	"github.com/zkMeLabs/mechain-relayer/types"
)

// VerifySignature verifies vote signature
func VerifySignature(vote *votepool.Vote, eventHash []byte) error {
	blsPubKey, err := bls.UnmarshalPublicKey(vote.PubKey[:])
	if err != nil {
		return errors.Wrap(err, "convert public key from bytes to bls failed")
	}
	sig, err := bls.UnmarshalSignature(vote.Signature[:])
	if err != nil {
		return errors.Wrap(err, "invalid signature")
	}
	if !sig.Verify(blsPubKey, eventHash[:], votepool.DST) {
		return errors.New("verify bls signature failed.")
	}
	return nil
}

// AggregateSignatureAndValidatorBitSet aggregates signature from multiple votes, and marks the bitset of validators who contribute votes
func AggregateSignatureAndValidatorBitSet(votes []*model.Vote, validators interface{}) ([]byte, *bitset.BitSet, error) {
	signatures := make(bls.Signatures, 0, len(votes))
	voteAddrSet := make(map[string]struct{}, len(votes))
	valBitSet := bitset.New(ValidatorsCapacity)
	for _, v := range votes {
		voteAddrSet[v.PubKey] = struct{}{}
		signature, _ := bls.UnmarshalSignature(common.Hex2Bytes(v.Signature))
		signatures = append(signatures, signature)
	}
	if reflect.TypeOf(validators).Elem() == reflect.TypeOf(types.Validator{}) {
		for idx, valInfo := range validators.([]types.Validator) {
			if _, ok := voteAddrSet[hex.EncodeToString(valInfo.BlsPublicKey[:])]; ok {
				valBitSet.Set(uint(idx))
			}
		}
	} else {
		for idx, valInfo := range validators.([]*tmtypes.Validator) {
			if _, ok := voteAddrSet[hex.EncodeToString(valInfo.BlsKey[:])]; ok {
				valBitSet.Set(uint(idx))
			}
		}
	}
	sigs, err := signatures.Aggregate().Marshal()
	if err != nil {
		return nil, valBitSet, err
	}
	return sigs, valBitSet, nil
}
