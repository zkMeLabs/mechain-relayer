package vote

import (
	"github.com/0xPolygon/polygon-edge/bls"
	"github.com/cometbft/cometbft/votepool"
)

type VoteSigner struct {
	privKey *bls.PrivateKey
	pubKey  *bls.PublicKey
}

// var DST = []byte("BLS_SIG_BN254G1_XMD:SHA-256_SVDW_RO_NUL_") //0x416a79f61c64a68f9946f715ec2b7077204a431ad6c623c2b7e464d4b60b0ed6

func NewVoteSigner(pk []byte) *VoteSigner {
	privKey, err := bls.UnmarshalPrivateKey(pk)
	if err != nil {
		panic(err)
	}
	pubKey := privKey.PublicKey()
	return &VoteSigner{
		privKey: privKey,
		pubKey:  pubKey,
	}
}

// SignVote signs a vote by relayer's private key
func (signer *VoteSigner) SignVote(vote *votepool.Vote) {
	vote.PubKey = signer.pubKey.Marshal()
	signature, err := signer.privKey.Sign(vote.EventHash[:], votepool.DST)
	if err != nil {
		panic(err)
	}
	vote.Signature, _ = signature.Marshal()
}
