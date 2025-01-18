package cmd

import (
	"os"

	"github.com/open-zhy/secm/internal/id"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate key",
	RunE:  runGenerate,
}

var (
	keyType string
	keySize int
)

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.PersistentFlags().StringVarP(&keyType, "type", "t", "rsa", "Key type, supports rsa, p256, p384, p521, ec25519")
	generateCmd.PersistentFlags().IntVar(&keySize, "size", 2048, "Key size, take effect for RSA key types only")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	identity, err := id.GenerateKey(
		id.GenerateKeyOpts{
			Type: keyType,
			Size: &keySize,
		},
	)

	if err != nil {
		return err
	}

	return identity.Encode(os.Stdout)
}
