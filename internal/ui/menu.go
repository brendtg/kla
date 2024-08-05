package ui

import (
	"strings"

	"github.com/manifoldco/promptui"
)

// DisplayMenu displays items using promptui for selection
func DisplayMenu(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
		Size:  50, // Increase the number of items to show in the list at a time
		Searcher: func(input string, index int) bool {
			item := items[index]
			return strings.Contains(strings.ToLower(item), strings.ToLower(input))
		},
	}

	_, result, err := prompt.Run()
	return result, err
}
