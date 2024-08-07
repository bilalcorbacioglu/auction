package cmd

import (
	"auction/x/auction/types"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func CmdCreateAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-auction [item] [starting-bid]",
		Short: "Create a new auction",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return fmt.Errorf("GetClientTxContext Error")
			}

			fromAddress := clientCtx.GetFromAddress().String()
			if fromAddress == "" {
				return fmt.Errorf("address cannot be empty")
			}

			item := args[0]
			startingBidStr := args[1]

			if item == "" {
				return fmt.Errorf("item cannot be empty")
			}

			startingBid, err := sdk.ParseCoinNormalized(startingBidStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateAuction(fromAddress, item, startingBid)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdPlaceBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "place-bid [auction-id] [bid-amount]",
		Short: "Place a bid on an auction",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return fmt.Errorf("GetClientTxContext Error")
			}

			fromAddress := clientCtx.GetFromAddress().String()
			if fromAddress == "" {
				return fmt.Errorf("address cannot be empty")
			}

			auctionID := args[0]
			bidAmountStr := args[1]

			bidAmount, err := sdk.ParseCoinNormalized(bidAmountStr)
			if err != nil {
				return err
			}

			msg := types.NewMsgPlaceBid(fromAddress, auctionID, bidAmount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
