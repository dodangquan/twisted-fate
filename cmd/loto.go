/*
Copyright © 2021 Đỗ Đăng Quân <dodangquan@outlook.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// lotoCmd represents the loto command
var lotoCmd = &cobra.Command{
	Use:   "loto",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cart := make(map[int64]interface{}, 0)

		for len(cart) < 6 {
			numTmp, err := rand.Int(rand.Reader, big.NewInt(maxNumberFlag))
			if err != nil {
				log.Fatal().Err(err).Send()
			}
			num := numTmp.Int64() + 1
			_, ok := cart[num]
			if !ok {
				cart[num] = true
			}
			time.Sleep(time.Duration(num*10) * time.Millisecond)
		}

		ticket := make([]int64, 0)
		for k, _ := range cart {
			ticket = append(ticket, k)
		}

		log.Info().Interface("ticket", ticket).Msg("Good luck! =))")
	},
}

var maxNumberFlag int64

func init() {
	rootCmd.AddCommand(lotoCmd)

	lotoCmd.Flags().Int64VarP(&maxNumberFlag, "max", "m",55, "")
}
