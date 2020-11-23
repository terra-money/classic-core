package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// MultiSignReq defines the properties of a multisign request's body.
type MultiSignReq struct {
	Tx            auth.StdTx          `json:"tx"`
	ChainID       string              `json:"chain_id"`
	Signatures    []auth.StdSignature `json:"signatures"`
	SignatureOnly bool                `json:"signature_only"`
	Sequence      uint64              `json:"sequence_number"`
	Pubkey        MultiSignPubKey     `json:"pubkey"` // (optional) In case the multisig account never reveals its pubkey, it is required.
}

// MultiSignPubKey defines the properties of a multisig account's public key
type MultiSignPubKey struct {
	Threshold int      `json:"threshold"`
	PubKeys   []string `json:"pubkeys"`
}

// MultiSignRequestHandlerFn - http request handler to build multisign transaction.
func MultiSignRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode params
		vars := mux.Vars(r)
		bech32Addr := vars["address"]
		accGetter := types.NewAccountRetriever(cliCtx)

		multiSignAddr, err := sdk.AccAddressFromBech32(bech32Addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := accGetter.EnsureExists(multiSignAddr); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		multiSignAccount, err := accGetter.GetAccount(multiSignAddr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Decode request body
		var req MultiSignReq
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(body, &req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var multisigPub multisig.PubKeyMultisigThreshold
		if multiSignAccount.GetPubKey() != nil {
			multisigPub = multiSignAccount.GetPubKey().(multisig.PubKeyMultisigThreshold)
		} else {

			var pubkeys []crypto.PubKey
			for _, bechPubkey := range req.Pubkey.PubKeys {
				pubkey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, bechPubkey)
				if err != nil {
					rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
					return
				}

				pubkeys = append(pubkeys, pubkey)
			}

			// Ensure threshold <= len(pubkeys)
			if req.Pubkey.Threshold > len(pubkeys) {
				err := fmt.Errorf("Not sufficient pubkeys; required: %d given: %d", req.Pubkey.Threshold, len(pubkeys))
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			multisigPub = multisig.NewPubKeyMultisigThreshold(req.Pubkey.Threshold, pubkeys).(multisig.PubKeyMultisigThreshold)
		}

		multisigSig := multisig.NewMultisig(len(multisigPub.PubKeys))

		accountNumber := multiSignAccount.GetAccountNumber()
		sequence := req.Sequence
		if req.Sequence == 0 {
			sequence = multiSignAccount.GetSequence()
		}

		if len(req.Signatures) < int(multisigPub.K) {
			err := fmt.Errorf("threashold: %v, # of given signatures: %v", len(req.Signatures), multisigPub.K)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// read each signature and add it to the multisig if valid
		for i := 0; i < len(req.Signatures); i++ {
			stdSig := req.Signatures[i]

			// Validate each signature
			sigBytes := auth.StdSignBytes(
				req.ChainID, accountNumber, sequence,
				req.Tx.Fee, req.Tx.GetMsgs(), req.Tx.GetMemo(),
			)
			if ok := stdSig.PubKey.VerifyBytes(sigBytes, stdSig.Signature); !ok {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "couldn't verify signature")
				return
			}
			if err := multisigSig.AddSignatureFromPubKey(stdSig.Signature, stdSig.PubKey, multisigPub.PubKeys); err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		newStdSig := auth.StdSignature{Signature: cliCtx.Codec.MustMarshalBinaryBare(multisigSig), PubKey: multisigPub}
		newTx := auth.NewStdTx(req.Tx.GetMsgs(), req.Tx.Fee, []auth.StdSignature{newStdSig}, req.Tx.GetMemo())

		sigOnly := req.SignatureOnly
		var json []byte
		switch {
		case sigOnly:
			json, err = cliCtx.Codec.MarshalJSONIndent(newTx.Signatures[0], "", "  ")
		case sigOnly && !cliCtx.Indent:
			json, err = cliCtx.Codec.MarshalJSON(newTx.Signatures[0])
		case !sigOnly && cliCtx.Indent:
			json, err = cliCtx.Codec.MarshalJSONIndent(newTx, "", "  ")
		default:
			json, err = cliCtx.Codec.MarshalJSON(newTx)
		}

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, json)
		return
	}
}
