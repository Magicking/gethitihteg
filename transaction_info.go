// Copyright 2017 Sylvain 6120 Laurent
// This file is part of the gethitihteg library.
//
// The gethitihteg library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The gethitihteg library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the gethitihteg library. If not, see <http://www.gnu.org/licenses/>.

package blockchain

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// Client defines typed wrappers for the Ethereum RPC API.
type Client struct {
	c *rpc.Client
}

// Dial connects a client to the given URL.
func Dial(rawurl string) (*Client, error) {
	c, err := rpc.Dial(rawurl)
	if err != nil {
		return nil, err
	}
	return newClient(c), nil
}

// NewClient creates a client that uses the given RPC client.
func newClient(c *rpc.Client) *Client {
	return &Client{c}
}

// HeaderByHash returns the block header with the given hash.
func (ec *Client) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	var head *types.Header
	err := ec.c.CallContext(ctx, &head, "eth_getBlockByHash", hash, false)
	if err == nil && head == nil {
		err = ethereum.NotFound
	}
	return head, err
}

// TransactionByHash returns the transaction with the given hash.
func (ec *Client) TransactionByHashFull(ctx context.Context, hash common.Hash) (tx *types.Transaction, blockhash *common.Hash, err error) {
	var raw json.RawMessage
	err = ec.c.CallContext(ctx, &raw, "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, nil, err
	} else if len(raw) == 0 {
		return nil, nil, ethereum.NotFound
	}
	if err := json.Unmarshal(raw, &tx); err != nil {
		return nil, nil, err
	} else if _, r, _ := tx.RawSignatureValues(); r == nil {
		return nil, nil, fmt.Errorf("server returned transaction without signature")
	}
	var block struct{ BlockHash *common.Hash }
	if err := json.Unmarshal(raw, &block); err != nil {
		return nil, nil, err
	}
	return tx, block.BlockHash, nil
}
