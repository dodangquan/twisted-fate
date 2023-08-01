/*
Package cmd
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
	"container/heap"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/dodangquan/twisted-fate/lucky"
	"github.com/pterm/pterm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	xrand "golang.org/x/exp/rand"
)

// luckyCmd represents the lucky command
var luckyCmd = &cobra.Command{
	Use:   "lucky",
	Short: "Lucky",
	Long:  `Lucky number`,
	Run: func(cmd *cobra.Command, args []string) {
		numberOfRotations := 6
		arr := make([]int64, maxNumberFlag)
		luckyNumbers := make([]string, 0)
		lockedSource := &xrand.LockedSource{}
		lockedSource.Seed(uint64(time.Now().UnixNano()))
		r := xrand.New(lockedSource)

		luckySpinner, err := pterm.DefaultSpinner.
			WithMessageStyle(pterm.NewStyle(pterm.FgDefault)).
			WithRemoveWhenDone(false).
			WithWriter(os.Stderr).
			Start("Finding lucky number...")
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		defer func() {
			e := luckySpinner.Stop()
			if e != nil {
				log.Fatal().Err(e).Send()
			}
		}()

		var idx int64 = 0
		for i := 0; i < numberOfRotations; i++ {
			pool := make(map[int64]interface{}, 0)
			idx = 0
			for idx < maxNumberFlag {
				randomize := r.Int63n(maxNumberFlag)
				num := randomize + 1
				_, ok := pool[num]
				if !ok {
					pool[num] = true
					arr[idx] = num
					idx++
				}
			}
		}

		h := createHeap(numberOfRotations, arr)

		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for i := 0; i < numberOfRotations; i++ {
			index := r.Int63n(int64(h.Len()))

			<-ticker.C
			num := heap.Remove(&h, int(index)).(*lucky.Number)

			luckyNumbers = append(luckyNumbers, fmt.Sprintf("%02d", num.Value))
			luckySpinner.Success(fmt.Sprintf("Number#%d: %02d", i+1, num.Value))

			arrTemp := make([]int64, h.Len())
			for j := 0; j < len(arrTemp); j++ {
				arrTemp[j] = h.Pop().(*lucky.Number).Value
			}

			h = createHeap(numberOfRotations, arrTemp)
		}

		sort.Strings(luckyNumbers)

		log.Info().Interface("LuckyNumbers", luckyNumbers).Msg("Good luck! =))")
		luckySpinner.Success("Finish")
	},
}

func createHeap(numberOfRotations int, arr []int64) lucky.NumberHeap {
	h := make(lucky.NumberHeap, numberOfRotations)
	for i := 0; i < numberOfRotations; i++ {
		h[i] = lucky.NewNumber(arr[i])
	}

	heap.Init(&h)
	for i := numberOfRotations; i < len(arr); i++ {
		heap.Push(&h, lucky.NewNumber(arr[i]))
	}

	return h
}

var maxNumberFlag int64

func init() {
	rootCmd.AddCommand(luckyCmd)

	luckyCmd.Flags().Int64VarP(&maxNumberFlag, "max", "m", 55, "")
}
