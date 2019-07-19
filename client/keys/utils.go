package keys

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
)

// available output formats.
const (
	OutputFormatText = "text"
	OutputFormatJSON = "json"

	// defaultKeyDBName is the client's subdirectory where keys are stored.
	defaultKeyDBName = "keys"
)

type bechKeyOutFn func(keyInfo keys.Info) (keys.KeyOutput, error)

// GetKeyInfo returns key info for a given name. An error is returned if the
// keybase cannot be retrieved or getting the info fails.
func GetKeyInfo(name string) (keys.Info, error) {
	keybase, err := NewKeyBaseFromHomeFlag()
	if err != nil {
		return nil, err
	}

	return keybase.Get(name)
}

// NewKeyBaseFromHomeFlag initializes a Keybase based on the configuration.
func NewKeyBaseFromHomeFlag() (keys.Keybase, error) {
	rootDir := viper.GetString(cli.HomeFlag)
	return NewKeyBaseFromDir(rootDir)
}

// NewKeyBaseFromDir initializes a keybase at a particular dir.
func NewKeyBaseFromDir(rootDir string) (keys.Keybase, error) {
	return getLazyKeyBaseFromDir(rootDir)
}

func getLazyKeyBaseFromDir(rootDir string) (keys.Keybase, error) {
	return keys.New(defaultKeyDBName, filepath.Join(rootDir, "keys")), nil
}

func printKeyTextHeader() {
	fmt.Printf("NAME:\tTYPE:\tADDRESS:\t\t\t\t\tPUBKEY:\n")
}

func printMultiSigKeyTextHeader() {
	fmt.Printf("WEIGHT:\tTHRESHOLD:\tADDRESS:\t\t\t\t\tPUBKEY:\n")
}

func printMultiSigKeyInfo(keyInfo keys.Info, bechKeyOut bechKeyOutFn) {
	ko, err := bechKeyOut(keyInfo)
	if err != nil {
		panic(err)
	}

	printMultiSigKeyTextHeader()
	printMultiSigKeyOutput(ko)
}

func printKeyInfo(keyInfo keys.Info, bechKeyOut bechKeyOutFn) {
	ko, err := bechKeyOut(keyInfo)
	if err != nil {
		panic(err)
	}

	switch viper.Get(cli.OutputFlag) {
	case OutputFormatText:
		printKeyTextHeader()
		printKeyOutput(ko)

	case OutputFormatJSON:
		var out []byte
		var err error
		if viper.GetBool(client.FlagIndentResponse) {
			out, err = cdc.MarshalJSONIndent(ko, "", "  ")
		} else {
			out, err = cdc.MarshalJSON(ko)
		}
		if err != nil {
			panic(err)
		}

		fmt.Println(string(out))
	}
}

func printInfos(infos []keys.Info) {
	kos, err := keys.Bech32KeysOutput(infos)
	if err != nil {
		panic(err)
	}

	switch viper.Get(cli.OutputFlag) {
	case OutputFormatText:
		printKeyTextHeader()
		for _, ko := range kos {
			printKeyOutput(ko)
		}

	case OutputFormatJSON:
		var out []byte
		var err error

		if viper.GetBool(client.FlagIndentResponse) {
			out, err = cdc.MarshalJSONIndent(kos, "", "  ")
		} else {
			out, err = cdc.MarshalJSON(kos)
		}

		if err != nil {
			panic(err)
		}
		fmt.Println(string(out))
	}
}

func printKeyOutput(ko keys.KeyOutput) {
	fmt.Printf("%s\t%s\t%s\t%s\n", ko.Name, ko.Type, ko.Address, ko.PubKey)
}

func printMultiSigKeyOutput(ko keys.KeyOutput) {
	for _, pk := range ko.PubKeys {
		fmt.Printf("%d\t%d\t\t%s\t%s\n", pk.Weight, ko.Threshold, pk.Address, pk.PubKey)
	}
}

func printKeyAddress(info keys.Info, bechKeyOut bechKeyOutFn) {
	ko, err := bechKeyOut(info)
	if err != nil {
		panic(err)
	}

	fmt.Println(ko.Address)
}

func printPubKey(info keys.Info, bechKeyOut bechKeyOutFn) {
	ko, err := bechKeyOut(info)
	if err != nil {
		panic(err)
	}

	fmt.Println(ko.PubKey)
}
