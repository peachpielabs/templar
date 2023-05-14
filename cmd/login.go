/*
Copyright (c) 2023 Peach Pie Labs, LLC.
*/

package cmd

// import (
// 	"errors"
// 	"fmt"
// 	"strconv"

// 	"github.com/manifoldco/promptui"
// 	"github.com/spf13/cobra"
// )

// func init() {
// 	rootCmd.AddCommand(loginCmd)
// }

// var loginCmd = &cobra.Command{
// 	Use:   "login",
// 	Short: "Log into GitFormer",
// 	Long:  `Log into GitFormer`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		validate := func(input string) error {
// 			_, err := strconv.ParseFloat(input, 64)
// 			if err != nil {
// 				return errors.New("Invalid number")
// 			}
// 			return nil
// 		}

// 		prompt := promptui.Prompt{
// 			Label:    "Number",
// 			Validate: validate,
// 		}

// 		result, err := prompt.Run()

// 		if err != nil {
// 			fmt.Printf("Prompt failed %v\n", err)
// 			return
// 		}

// 		fmt.Printf("You choose %q\n", result)
// 	},
// }
