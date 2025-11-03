package cmd

import (
	"algo/internal/model"
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"os"
	"strings"
)

func getCmdParam(cmd *cobra.Command, name, prompt string, debug bool) string {
	value := getOrPrompt(cmd.Flag(name).Value.String(), prompt)
	if debug {
		fmt.Printf("%s: [%s]\n", name, value)
	}
	return value
}

func getOrPrompt(flagValue, prompt string) string {
	if strings.TrimSpace(flagValue) != "" {
		return strings.TrimSpace(flagValue)
	}

	fmt.Print(prompt)
	return Scanln()
}

func Scanln() string {
	reader := bufio.NewReader(os.Stdin)
	var lines []string

	for {
		line, _ := reader.ReadString('\n')
		line = strings.TrimRight(line, "\r\n")

		if strings.HasSuffix(line, "\\") {
			lines = append(lines, line[:len(line)-1])
			fmt.Print("> ") // 提示用户输入续行
			continue
		}

		lines = append(lines, line)
		break
	}

	return strings.Join(lines, "")
}

func clearProblemTags(tx *gorm.DB, problem *model.Problem) error {
	// 先预加载 Tags
	if err := tx.Preload("Tags").First(problem, problem.ID).Error; err != nil {
		return fmt.Errorf("failed to load problem tags: %w", err)
	}

	// 清空多对多关联
	if err := tx.Model(problem).Association("Tags").Clear(); err != nil {
		return fmt.Errorf("failed to clear tags association: %w", err)
	}
	return nil
}
