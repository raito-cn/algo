package main

import (
	"algo/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "algo",
		Short: "[ algo是一个用于管理算法题的命令行工具 ] algo is a command line tool for algorithm",
	}

	rootCmd.AddCommand(cmd.InitAddCmd())
	rootCmd.AddCommand(cmd.InitListCmd())
	rootCmd.AddCommand(cmd.InitRemoveCmd())
	rootCmd.AddCommand(cmd.InitEditCmd())
	rootCmd.AddCommand(cmd.InitGenCmd())
	_ = rootCmd.Execute()
}
