package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"

	"github.com/yaitoo/sparrow/db/util"

	"github.com/manifoldco/promptui"
)

// var secretLength = flag.Int("length", 0, "Input Secret Length")
var secretContent = flag.String("secret", "", "Input Your Secret")
var plaintText = flag.String("plain", "", "Input Your Plaintext")

func main() {
	flag.Parse()

	chaotext, err := util.AesEncrypt(*plaintText, *secretContent)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(chaotext)

	prompt := promptui.Select{
		Label: "Select Secret Length",
		Items: []int{16, 24, 32},
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "\U0001F336 {{ . | cyan }} ",
			Inactive: "\U0001F336 {{ . |  yellow}} ",
		},
	}

	_, secretLength, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	fmt.Printf("You choose length is %q\n", secretLength)

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	promptInputSecret := promptui.Prompt{
		Label:     "Input Secret",
		Templates: templates,
		AllowEdit: true,
		Validate: func(input string) error {
			length, _ := strconv.Atoi(secretLength)
			if len(input) == length {
				return nil
			}
			return errors.New("secret must have " + secretLength + " characters")
		},
	}

	inputSecret, _ := promptInputSecret.Run()
	fmt.Printf("You input %q\n", inputSecret)

	promptInputPlaintext := promptui.Prompt{
		Label:     "Input Plaintext",
		Templates: templates,
		AllowEdit: true,
	}
	inputPlaintext, _ := promptInputPlaintext.Run()
	fmt.Printf("You input %q\n", inputPlaintext)
	result, err := util.AesEncrypt(inputPlaintext, inputSecret)
	fmt.Printf("You cipher %q\n", result)
}
