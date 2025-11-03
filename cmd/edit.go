package cmd

import (
	"algo/internal/db"
	"algo/internal/model"
	"fmt"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"strconv"
)

var editCmd = &cobra.Command{
	Use:   "edit [slug]",
	Short: "[ 编辑题目 ] Edit a problem",
	Args:  cobra.ExactArgs(1),
	Run:   editProblem,
}

func InitEditCmd() *cobra.Command {
	editCmd.Long = `Edit a problem by its slug.
Example:
  algo edit two-sum --debug`

	editCmd.Flags().StringP("title", "t", "", "[ 题目标题 ] Problem title")
	editCmd.Flags().StringP("difficulty", "d", "", "[ 题目难度 ] Problem difficulty (easy|medium|hard)")
	editCmd.Flags().StringP("tags", "g", "", "[ 题目标签，英文逗号分割 ] Problem tags, comma separation")
	editCmd.Flags().StringP("solution", "S", "", "[ 题目在线地址 ] Problem solution URL")
	editCmd.Flags().StringP("note", "n", "", "[ 题目笔记 ] Problem note")
	editCmd.Flags().StringP("codePath", "c", "", "[ 题目代码本地地址 ] Problem code path")
	editCmd.Flags().BoolP("debug", "D", false, "Debug mode")
	editCmd.Flags().StringP("score", "s", "", "[ 题目评分 ] Problem score")

	return editCmd
}

func editProblem(cmd *cobra.Command, args []string) {
	debug, _ := cmd.Flags().GetBool("debug")

	title := cmd.Flag("title").Value.String()
	difficulty := cmd.Flag("difficulty").Value.String()
	tags := cmd.Flag("tags").Value.String()
	solution := cmd.Flag("solution").Value.String()
	note := cmd.Flag("note").Value.String()
	codePath := cmd.Flag("codePath").Value.String()
	score := cmd.Flag("score").Value.String()

	// 获取 slug
	slug := args[0]

	conn := db.GetDB(debug)
	_ = conn.Transaction(func(tx *gorm.DB) error {
		var problem model.Problem
		if err := tx.Where("slug = ?", slug).First(&problem).Error; err != nil {
			return fmt.Errorf("problem not found: %w", err)
		}

		if title != "" {
			problem.Title = title
			problem.SetSlug()
		}
		if difficulty != "" {
			diff := model.Difficulty(difficulty)
			if !diff.Valid() {
				return fmt.Errorf("invalid difficulty")
			}
			problem.Difficulty = diff
		}
		if tags != "" {
			checkTags, err := CheckTags(tx, tags)
			if err != nil {
				return err
			}
			if err = clearProblemTags(tx, &problem); err != nil {
				return fmt.Errorf("failed to clear old tags: %w", err)
			}
			// 重新添加新标签
			if err = tx.Model(&problem).Association("Tags").Replace(checkTags); err != nil {
				return fmt.Errorf("failed to update tags: %w", err)
			}
			problem.Tags = checkTags
		}
		if solution != "" {
			problem.SolutionURL = solution
		}
		if note != "" {
			problem.Note = note
		}
		if codePath != "" {
			problem.CodePath = codePath
			problem.SetCodePath()
		}
		if score != "" {
			atoi, err := strconv.Atoi(score)
			if err != nil {
				return err
			}
			if atoi >= 0 && atoi <= 255 {
				v := uint8(atoi)
				problem.Score = &v
			}
		}

		if err := tx.Save(&problem).Error; err != nil {
			return fmt.Errorf("failed to edit problem: %w", err)
		}

		fmt.Println("Problem edited successfully")
		return nil
	})
}
