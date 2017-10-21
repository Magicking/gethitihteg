// Copyright 2017 Sylvain 6120 Laurent
// This file is part of the gethitihteg library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/jessevdk/go-flags"

	ethtk "github.com/Magicking/gethitihteg"
)

var opts struct {
	RpcURL     string   `long:"rpc-url" default:"http://localhost:8545" description:"RPC URL for the node"`
	PrivateKey string   `long:"key" required:"true" description:"Private key used to sign transaction"`
	SolF       []string `long:"sol" required:"true" description:"Path to the Ethereum contract Solidity source to deploy"`
	SolC       string   `long:"solc" default:"solc" description:"Solidity compiler to use if source builds are requested"`
	Excl       []string `long:"excl" description:"Contract to exclude from deploying"`
}

func ParamsParse(params []string) (ret []interface{}, err error) {
	//TODO ret[contract_name][]interface{}
	for _, e := range params {
		val := strings.SplitN(e, ":", 2)
		if len(val) != 2 {
			return nil, fmt.Errorf("Bad argument format for %q should be \"type:value\"", e)
		}
		switch val[0] {
		case "string":
			buf := val[1]
			ret = append(ret, &buf)
		case "int":
			var v int64
			fmt.Sscan(val[1], &v)
			ret = append(ret, big.NewInt(v))
			//TODO types
		default:
			return nil, fmt.Errorf("Unknown type %q", val[0])
		}
	}
	return ret, nil
}

func main() {
	lst, err := flags.Parse(&opts)
	if err != nil {
		return
	}
	client, err := ethtk.NewNodeConnector(opts.RpcURL, 3)
	if err != nil {
		log.Fatalf("Could not initialize client context: %v", err)
	}
	params, err := ParamsParse(lst)
	if err != nil {
		log.Fatal(err)
	}
	//Deploy contract
	contracts, err := compiler.CompileSolidity(opts.SolC, opts.SolF...)
	if err != nil {
		log.Fatalf("Failed to build Solidity contract: %v\n", err)
	}
	excludedContracts := make(map[string]bool)
	for _, e := range opts.Excl {
		excludedContracts[strings.ToLower(e)] = true
	}
	// Gather all non-excluded contract for binding
	var TODO bool
	for name, contract := range contracts {
		if TODO {
			log.Fatal("TODO handle only one contract")
		}
		if excludedContracts[strings.ToLower(name)] {
			continue
		}
		abi, _ := json.Marshal(contract.Info.AbiDefinition) // Flatten the compiler parse
		nameParts := strings.Split(name, ":")
		fmt.Println(nameParts)
		addr, tx, _, err := ethtk.CreateContractHelper(client, opts.PrivateKey, string(abi), contract.Code, params...)
		fmt.Println(addr, tx, err)

		TODO = true
	}
	fmt.Println(opts)
}
