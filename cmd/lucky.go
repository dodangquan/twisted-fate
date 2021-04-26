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
	"fmt"
	"math/big"
	"time"

	"github.com/pterm/pterm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// luckyCmd represents the lucky command
var luckyCmd = &cobra.Command{
	Use:   "lucky",
	Short: "Lucky",
	Long:  `Lucky number`,
	Run: func(cmd *cobra.Command, args []string) {
		cart := make(map[int64]interface{}, 0)
		arr := make([]int64, maxNumberFlag)
		ticket := make([]string, 0)
		var idx int64 = 0

		luckySpinner, err := pterm.DefaultSpinner.WithMessageStyle(pterm.NewStyle(pterm.FgDefault)).Start("Finding lucky number...")
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		defer func() {
			e := luckySpinner.Stop()
			if e != nil {
				log.Fatal().Err(e).Send()
			}
		}()

		for i := 0; i < 6; i++ {
			pool := make(map[int64]interface{}, 0)
			idx = 0
			for ; idx < maxNumberFlag; {
				numTmp, err := rand.Int(rand.Reader, big.NewInt(maxNumberFlag))
				if err != nil {
					log.Fatal().Err(err).Send()
				}
				num := numTmp.Int64() + 1
				_, ok := pool[num]
				if !ok {
					time.Sleep(time.Duration(num*20) * time.Millisecond)
					pool[num] = true
					arr[idx] = num
					idx++
				}
			}
		}

		for len(cart) < 6 {
			numTmp, err := rand.Int(rand.Reader, big.NewInt(maxNumberFlag))
			if err != nil {
				log.Fatal().Err(err).Send()
			}
			num := numTmp.Int64()
			_, ok := cart[num]
			if !ok {
				cart[num] = true
				luckySpinner.Success(fmt.Sprintf("Number#%d: %02d", len(cart), arr[num]))
				time.Sleep(time.Duration(num*50) * time.Millisecond)
			}
		}
		for k := range cart {
			ticket = append(ticket, fmt.Sprintf("%02d", arr[k]))
		}

		log.Info().Interface("ticket", ticket).Msg("Good luck! =))")
		luckySpinner.Success("Finish")
	},
}

var maxNumberFlag int64

func init() {
	rootCmd.AddCommand(luckyCmd)

	luckyCmd.Flags().Int64VarP(&maxNumberFlag, "max", "m", 55, "")
}
