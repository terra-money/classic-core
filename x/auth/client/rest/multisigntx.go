package cli

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/multisig"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, kb keys.Keybase) {
	r.HandleFunc("/auth/accounts/{address}/multisign", MultiSignRequestHandlerFn(cdc, kb, cliCtx)).Methods("POST")
}

// MultiSignReq defines the properties of a multisign request's body.
type MultiSignReq struct {
	Tx            auth.StdTx          `json:"tx"`
	ChainID       string              `json:"chain_id"`
	Signatures    []auth.StdSignature `json:"signatures"`
	SignatureOnly bool                `json:"signature_only"`
	Sequence      uint64              `json:"sequence_number"`
	Pubkey        string              `json:"pubkey"` // (optional) In case the multisig account never reveals its pubkey, it is required.
}

// MultiSignRequestHandlerFn - http request handler to build multisign transaction.
func MultiSignRequestHandlerFn(cdc *codec.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode params
		vars := mux.Vars(r)
		bech32Addr := vars["address"]

		multiSignAddr, err := sdk.AccAddressFromBech32(bech32Addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		multiSignAccount, err := cliCtx.GetAccount(multiSignAddr.Bytes())
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

		err = cdc.UnmarshalJSON(body, &req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var multisigPub multisig.PubKeyMultisigThreshold
		if multiSignAccount.GetPubKey() != nil {
			multisigPub = multiSignAccount.GetPubKey().(multisig.PubKeyMultisigThreshold)
		} else {
			pubKey, err := sdk.GetAccPubKeyBech32(req.Pubkey)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			multisigPub = pubKey.(multisig.PubKeyMultisigThreshold)
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

		newStdSig := auth.StdSignature{Signature: cdc.MustMarshalBinaryBare(multisigSig), PubKey: multisigPub}
		newTx := auth.NewStdTx(req.Tx.GetMsgs(), req.Tx.Fee, []auth.StdSignature{newStdSig}, req.Tx.GetMemo())

		sigOnly := req.SignatureOnly
		var json []byte
		switch {
		case sigOnly:
			json, err = cdc.MarshalJSONIndent(newTx.Signatures[0], "", "  ")
		case sigOnly && !cliCtx.Indent:
			json, err = cdc.MarshalJSON(newTx.Signatures[0])
		case !sigOnly && cliCtx.Indent:
			json, err = cdc.MarshalJSONIndent(newTx, "", "  ")
		default:
			json, err = cdc.MarshalJSON(newTx)
		}

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, json, cliCtx.Indent)
		return
	}
}

func readAndUnmarshalStdSignature(cdc *amino.Codec, filename string) (stdSig auth.StdSignature, err error) {
	var bytes []byte
	if bytes, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	if err = cdc.UnmarshalJSON(bytes, &stdSig); err != nil {
		return
	}
	return
}
