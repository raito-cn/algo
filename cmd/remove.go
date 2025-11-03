package cmd

import (
	"algo/internal/db"
	"algo/internal/model"
	"fmt"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"os"
)

var removeCmd = &cobra.Command{
	Use:   "rm [slug]",
	Short: "[ 删除题目 ] Remove a problem",
	Args:  cobra.ExactArgs(1),
	Run:   deleteProblem,
}

func InitRemoveCmd() *cobra.Command {
	removeCmd.Flags().BoolP("debug", "D", false, "Debug mode")
	removeCmd.Long = `Remove a problem by its slug.
Example:
  algo rm 0001_two-sum --debug`
	return removeCmd
}

func deleteProblem(cmd *cobra.Command, args []string) {
	debug, _ := cmd.Flags().GetBool("debug")
	slug := args[0]

	conn := db.GetDB(debug)

	err := conn.Transaction(func(tx *gorm.DB) error {
		var problem model.Problem
		if err := tx.Where("slug = ?", slug).First(&problem).Error; err != nil {
			return fmt.Errorf("problem not found: %w", err)
		}

		// 清空标签关联
		if err := clearProblemTags(tx, &problem); err != nil {
			return err
		}

		// 删除题目
		if err := tx.Delete(&problem).Error; err != nil {
			return fmt.Errorf("failed to delete problem: %w", err)
		}

		if problem.CodePath != "" {
			if err := os.Remove(problem.CodePath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove code file: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Failed to remove problem:", err)
		return
	}

	fmt.Println("Problem removed successfully")
}
