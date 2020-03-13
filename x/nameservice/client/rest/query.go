package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/terra-project/core/x/nameservice/internal/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/nameservice/names/{%s}/auction", RestName), queryAuctionHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/nameservice/names/{%s}/auction/bids", RestName), queryBidHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/nameservice/names/{%s}/auction/bids/{%s}", RestName, RestBidderAddr), queryBidHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/nameservice/names/{%s}/registry", RestName), queryRegistryHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/nameservice/names/{%s}/resolve", RestName), queryResolveHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/nameservice/auctions", queryAuctionHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/nameservice/auctions/{%s}", RestStatus), queryAuctionHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/nameservice/addresses/{%s}/registry", RestAddress), queryReverseHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/nameservice/parameters", queryParamsHandlerFn(cliCtx)).Methods("GET")
}

func queryAuctionHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		nameStr := vars[RestName]

		var bz []byte
		if len(nameStr) != 0 {
			name := types.Name(nameStr)
			if name.Levels() != 2 {
				err := fmt.Errorf("must submit by the second level name")
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			nameHash, _ := name.NameHash()
			params := types.QueryAuctionsParams{NameHash: nameHash}
			bz = cliCtx.Codec.MustMarshalJSON(params)
		} else {
			status := types.AuctionStatusNil
			statusStr := vars[RestStatus]
			if len(statusStr) != 0 {
				var err error
				status, err = types.AuctionStatusFromString(statusStr)
				if err != nil {
					rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
					return
				}
			}

			params := types.QueryAuctionsParams{Status: status}
			bz = cliCtx.Codec.MustMarshalJSON(params)
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAuctions), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryBidHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		nameStr := vars[RestName]
		bidderAddrStr := vars[RestBidderAddr]

		name := types.Name(nameStr)
		if name.Levels() != 2 {
			err := fmt.Errorf("must submit by the second level name")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var bz []byte

		if len(bidderAddrStr) == 0 {
			nameHash, _ := name.NameHash()
			params := types.QueryBidsParams{NameHash: nameHash, Bidder: nil}
			bz = cliCtx.Codec.MustMarshalJSON(params)
		} else {
			bidderAddr, err := sdk.AccAddressFromBech32(bidderAddrStr)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			nameHash, _ := name.NameHash()
			params := types.QueryBidsParams{NameHash: nameHash, Bidder: bidderAddr}
			bz = cliCtx.Codec.MustMarshalJSON(params)
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryBids), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryRegistryHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		nameStr := vars[RestName]

		name := types.Name(nameStr)
		if name.Levels() != 2 {
			err := fmt.Errorf("must submit by the second level name")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		nameHash, _ := name.NameHash()
		params := types.QueryRegistryParams{NameHash: nameHash}
		bz := cliCtx.Codec.MustMarshalJSON(params)

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRegistry), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryResolveHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		nameStr := vars[RestName]

		name := types.Name(nameStr)
		if levels := name.Levels(); levels != 2 && levels != 3 {
			err := fmt.Errorf("must submit by the second or third level name")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		nameHash, childNameHash := name.NameHash()
		params := types.QueryResolveParams{NameHash: nameHash, ChildNameHash: childNameHash}
		bz := cliCtx.Codec.MustMarshalJSON(params)

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryResolve), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryReverseHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		addrStr := vars[RestAddress]
		addr, err := sdk.AccAddressFromBech32(addrStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		params := types.QueryReverseParams{Address: addr}
		bz := cliCtx.Codec.MustMarshalJSON(params)

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryReverse), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryParamsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
