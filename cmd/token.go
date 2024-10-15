package cmd

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sevensolutions/tiny-repo/core"
	"github.com/spf13/cobra"
)

var namespace string

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Interact with access tokens",
	Long:  `Interact with access tokens`,
}

var tokenCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new access token",
	Long:  `Create a new access token`,
	Run: func(cmd *cobra.Command, args []string) {

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"namespace": namespace,
			"name":      "Unknown",
			"prefix":    "/",
			"iat":       time.Now().Unix(),
		})

		secret := []byte(core.GetRequiredEnvVar("JWT_SECRET"))

		token, err := jwtToken.SignedString(secret)
		if err != nil {
			panic(err)
		}

		println(token)
	},
}

var tokenInspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect an access token",
	Long:  `Inspect an access token`,
	Run: func(cmd *cobra.Command, args []string) {

		tokenString := args[0]

		secret := []byte(core.GetRequiredEnvVar("JWT_SECRET"))

		claims := jwt.MapClaims{}
		jwtToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil {
			panic(err)
		}

		if jwtToken.Valid {
			println("The token is valid!")
		} else {
			println("The token is invalid!")
		}

		println("Claims:")
		for key, val := range claims {
			fmt.Printf("%v: %v\n", key, val)
		}
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)

	tokenCreateCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "The namespace the token should have access to")

	tokenCmd.AddCommand(tokenCreateCmd)
	tokenCmd.AddCommand(tokenInspectCmd)
}
