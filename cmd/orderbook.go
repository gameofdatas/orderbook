package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"pricer/models"
	"pricer/pricer"
)

// Command line Args
var runCmd = &cobra.Command{
	Use:   "pricer",                 // SubCommand
	Short: "Build limit order book", // Short description of the SubCommand
	Long:  "Limit order book implementation",
	RunE:  Pricer,
}

// inputPath for file of Order book events.
var inputPath string

// targetSize is total size shares to be sold or bought for a side
var targetSize int64

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&inputPath, "i", "i", "", "input path of the file")
	err := runCmd.MarkFlagRequired("i")
	handleCommandError(err)
	runCmd.Flags().Int64VarP(&targetSize, "target", "t", 0, "target size shares to be sold/bought")
	err = runCmd.MarkFlagRequired("target")
	handleCommandError(err)
}

// Pricer process the orders coming in
func Pricer(_ *cobra.Command, _ []string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open the pricer input file")
	}
	defer file.Close()
	buyObject := pricer.NewPricer(models.Buy, targetSize)
	sellObject := pricer.NewPricer(models.Sell, targetSize)
	reader := bufio.NewScanner(file)
	for reader.Scan() {
		err = handleOrderEvents(reader.Text(), buyObject, sellObject)
		if err != nil {
			return err
		}
	}
	if err := reader.Err(); err != nil {
		return err
	}
	return nil
}
