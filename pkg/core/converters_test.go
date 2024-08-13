package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/tonkeeper/opentonapi/internal/g"
	"github.com/tonkeeper/tongo/abi"
	"github.com/tonkeeper/tongo/boc"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/tlb"
)

func readFile[T any](filename string) (*T, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var t T
	if err := json.Unmarshal(bytes, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func TestConvertToAccount(t *testing.T) {
	tests := []struct {
		name            string
		accountID       tongo.AccountID
		filename        string
		want            *Account
		wantCodePresent bool
		wantDataPresent bool
	}{
		{
			name:      "active account with data",
			filename:  "testdata/account.json",
			accountID: tongo.MustParseAccountID("EQDendoireMDFMufOUzkqNpFIay83GnjV2tgGMbA64wA3siV"),
			want: &Account{
				AccountAddress:    tongo.MustParseAccountID("EQDendoireMDFMufOUzkqNpFIay83GnjV2tgGMbA64wA3siV"),
				Status:            tlb.AccountActive,
				TonBalance:        989109352,
				ExtraBalances:     nil,
				LastTransactionLt: 31236013000006,
				Storage: StorageInfo{
					UsedCells:       *big.NewInt(46),
					UsedBits:        *big.NewInt(13485),
					UsedPublicCells: *big.NewInt(0),
					LastPaid:        1663270333,
					DuePayment:      0,
				},
			},
			wantDataPresent: true,
			wantCodePresent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := readFile[tlb.ShardAccount](tt.filename)
			require.Nil(t, err)
			got, err := ConvertToAccount(tt.accountID, *account)
			require.Nil(t, err)
			if tt.wantCodePresent {
				require.True(t, len(got.Code) > 0)
			} else {
				require.Nil(t, got)
			}
			if tt.wantDataPresent {
				require.True(t, len(got.Data) > 0)
			} else {
				require.Nil(t, got)
			}
			got.Code = nil
			got.Data = nil
			require.Equal(t, tt.want, got)
		})
	}
}

func TestConvertTransaction(t *testing.T) {
	tests := []struct {
		name           string
		accountID      tongo.AccountID
		txHash         tongo.Bits256
		txLt           uint64
		filenamePrefix string
		wantErr        bool
	}{
		{
			name:           "convert-tx-1",
			txHash:         tongo.MustParseHash("6c41096bbe0c2ca57f652ca7362a43473f8b33d8fa555a673bc70bb85fab37f6"),
			txLt:           37362820000003,
			accountID:      tongo.MustParseAccountID("0:6dcb8357c6bef52b43f0f681d976f5a46068ae195cb95f7a959d25c71b0cac6c"),
			filenamePrefix: "convert-tx-1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, err := liteapi.NewClient(liteapi.FromEnvsOrMainnet())
			require.Nil(t, err)
			txs, err := cli.GetTransactions(context.Background(), 1, tt.accountID, tt.txLt, tt.txHash)
			require.Nil(t, err)
			tx, err := ConvertTransaction(0, txs[0])
			require.Nil(t, err)
			bs, err := json.MarshalIndent(tx, " ", "  ")
			require.Nil(t, err)
			outputName := fmt.Sprintf("testdata/%v.output.json", tt.filenamePrefix)
			if err := os.WriteFile(outputName, bs, 0644); err != nil {
				t.Fatalf("os.WriteFile() failed: %v", err)
			}
			expected, err := os.ReadFile(fmt.Sprintf("testdata/%v.json", tt.filenamePrefix))
			require.Nil(t, err)
			if bytes.Compare(bs, expected) != 0 {
				t.Fatalf("dont match")
			}
		})
	}
}

func TestConvertMessage(t *testing.T) {
	tests := []struct {
		name    string
		msgHash tongo.Bits256
		rawBody string
		wantOp  *abi.MsgOpCode
	}{
		{
			name:    "WalletSignedV4",
			msgHash: tongo.MustParseHash("fcd573be1b5b9212fbf296958c2ee7913c64175169b1ba58679ae69f60d69399"),
			rawBody: "b5ee9c720101020100a100019c2f4f2a62f97346c45e5232d66c8d525deaed365ab53f0c52753a370b0224f5f1e022c77b3d4ccfbdd0b007aedc5dc493e7ac98283cc6b47ca01c977a92764c0729a9a31766b30c8d00000073000301009c620008b38c316ada97ff0e0c2d8d3fc27aacb654403401255e94a57f97b56be136909138800000000000000000000000000000000000363662333063343864653535393336633533303532333134",
		},
		{
			name:    "WalletSignedExternalV5R1",
			msgHash: tongo.MustParseHash("66ea3850b7f188559ca6e657035183b3d5f402e3d0fc8e6cd478e7d10cc797d9"),
			rawBody: "b5ee9c720101040100950001a17369676e7fffff11ffffffff00000000a151e8192f8651c72c4825350f18e31510044897899b761dc2e9606283f358e58b742d7c8e6baf4ca68d0e2d5ad7f62d80688f166ec0c720a5ac59efa1041143a001020a0ec3c86d0302030000006862005304a75fe829b8489154323ea3f1577055dcce862eb7aa8cbec7e892909fc6a3a02faf080000000000000000000000000000",
			wantOp:  g.Pointer(uint32(0x7369676e)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roots, err := boc.DeserializeBocHex(tt.rawBody)
			require.Nil(t, err)
			if len(roots) != 1 {
				t.Fatalf("invaild raw body")
			}
			root := roots[0]

			var decodedBody *DecodedMessageBody
			tag, op, value, err := abi.ExtInMessageDecoder(root, extInMsgDecoderInterfaces)
			if err != nil || op == nil {
				t.Fatalf("decode message failed: %v", err)
			}
			decodedBody = &DecodedMessageBody{
				Operation: *op,
				Value:     value,
			}

			require.Equal(t, tt.wantOp, tag)
			require.Equal(t, tt.name, decodedBody.Operation)
		})
	}
}
