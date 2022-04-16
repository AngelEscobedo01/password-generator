package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/rivo/tview"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

// define all the global variables
var (
	// lowerCaseList               = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
	// upperCharList               = []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
	lowerCharSet                = "abcdefghijklmnopqrstuvwxyz"
	upperCharSet                = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet              = []rune{'!', '@', '#', '$', '%', '^', '&', '*', '(', ')', ',', '-'}
	numberCharSet               = "123567890"
	numberCharList              = []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9'}
	minSpecialChar              = 2
	minUpperChar                = 2
	minNumberChar               = 2
	passwordLength              = 10
	customSpecialCharacterTypes = ""
	customLowerCaseTypes        = ""
	customUpperCaseTypes        = ""
)

func main() {

	app := tview.NewApplication()
	pages := tview.NewPages()

	initPages(pages, app)

	//check if the password length matches the criteria
	totalCharLenWithoutLowerChar := minUpperChar + minSpecialChar + minNumberChar

	if totalCharLenWithoutLowerChar >= passwordLength {
		fmt.Println("Please provide valid password length")
		os.Exit(1)
	}

	// Get the user input - target folder needs to be organized
	scanner := bufio.NewScanner(os.Stdin)
	// fmt.Printf("How many passwords you want to generate? - ")
	// scanner.Scan()

	numberOfPasswords, err := strconv.Atoi(scanner.Text())

	if err != nil {
		// fmt.Println("Please provide correct value for number of passwords")
		os.Exit(1)
	}

	// it generate random number e
	rand.Seed(time.Now().Unix())

	for i := 0; i < numberOfPasswords; i++ {
		password := generatePassword()
		entropy := passwordvalidator.GetEntropy(password)
		fmt.Printf("Password %v is %v \n", i+1, password)
		fmt.Printf("Strength of password is %v \n", entropy)
	}

}

func generatePassword() string {

	// declare empty password variable
	password := ""

	// generate random special character based on minSpecialChar

	for i := 0; i < minSpecialChar; i++ {
		random := rand.Intn(len(specialCharSet))
		//fmt.Println(specialCharSet[random])
		//fmt.Printf("%v and %T \n", random, specialCharSet[random])
		password = password + string(specialCharSet[random])
	}

	// generate random upper character based on minUpperChar
	for i := 0; i < minUpperChar; i++ {
		random := rand.Intn(len(upperCharSet))
		password = password + string(upperCharSet[random])
	}

	// generate random upper character based on minNumberChar
	for i := 0; i < minNumberChar; i++ {
		random := rand.Intn(len(numberCharSet))
		password = password + string(numberCharSet[random])
	}

	// find remaining lowerChar
	totalCharLenWithoutLowerChar := minUpperChar + minSpecialChar + minNumberChar

	remainingCharLen := passwordLength - totalCharLenWithoutLowerChar

	// generate random lower character based on remainingCharLen
	for i := 0; i < remainingCharLen; i++ {
		random := rand.Intn(len(lowerCharSet))
		password = password + string(lowerCharSet[random])
	}

	// shuffle the password string

	passwordRune := []rune(password)
	rand.Shuffle(len(passwordRune), func(i, j int) {
		passwordRune[i], passwordRune[j] = passwordRune[j], passwordRune[i]
	})

	password = string(passwordRune)
	return password
}

func initPages(pages *tview.Pages, app *tview.Application) {

	customPasswordForm := tview.NewForm().
		AddInputField("Types of special characters", "", 30,
			func(textToCheck string, lastChar rune) bool {
				v, _ := strconv.Atoi(textToCheck)
				if v > 12 {
					return false
				}
				return isIncludedInList(specialCharSet, lastChar)
			}, func(text string) {
				// Store special character types.
				customSpecialCharacterTypes = text
			}).
		AddInputField("Amount of lowercase characters (0-12)", "", 30,
			func(textToCheck string, lastChar rune) bool {
				v, _ := strconv.Atoi(textToCheck)
				if v > 12 {
					return false
				}
				return isIncludedInList(numberCharList, lastChar)
			},
			func(text string) {
				customLowerCaseTypes = text
			}).
		AddInputField("Amount of UPPERCASE characters (0-12)", "", 30,
			func(textToCheck string, lastChar rune) bool {
				v, _ := strconv.Atoi(textToCheck)
				if v > 12 {
					return false
				}
				return isIncludedInList(numberCharList, lastChar)
			},
			func(text string) {
				minUpperChar, _ = strconv.Atoi(text)
			}).
		AddInputField("Amount of numbers characters (0-12)", "", 30,
			func(textToCheck string, lastChar rune) bool {
				v, _ := strconv.Atoi(textToCheck)
				if v > 12 {
					return false
				}
				return isIncludedInList(numberCharList, lastChar)
			},
			func(text string) {
				minNumberChar, _ = strconv.Atoi(text)
			}).
		AddInputField("Length of password (0-48)", "", 30,
			func(textToCheck string, lastChar rune) bool {
				v, _ := strconv.Atoi(textToCheck)
				if v > 48 {
					return false
				}
				return isIncludedInList(numberCharList, lastChar)
			}, func(text string) {
				passwordLength, _ = strconv.Atoi(text)
			}).
		AddButton("Save", func() {
			app.Stop()
			password := generatePassword()
			entropy := passwordvalidator.GetEntropy(password)

			fmt.Printf("\n")
			fmt.Printf("Password is %v \n", password)
			fmt.Printf("Strength of password is %v \n", math.Round(entropy))
		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	pages.AddPage("page-0",
		tview.NewList().
			AddItem("Create default password?", "", 'a', func() {
				app.Stop()
				password := generatePassword()
				entropy := passwordvalidator.GetEntropy(password)

				fmt.Printf("\n")
				fmt.Printf("Password is %v \n", password)
				fmt.Printf("Strength of password is %v \n", math.Round(entropy))
			}).
			AddItem("Generate custom password?", "", 'b', func() {
				pages.SwitchToPage("page-2")
			}).
			AddItem("Exit", "", 'q', func() {
				app.Stop()
			}),
		true, true)
	pages.AddPage("page-1",
		tview.NewModal().
			SetText("This is page 1. Choose where to go next").
			AddButtons([]string{"Next", "Quit"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonIndex == 0 {
					pages.SwitchToPage("test")
				} else {
					app.Stop()
				}
			}),
		false, false)
	pages.AddPage("page-2", customPasswordForm, true, false)

	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}

func isIncludedInList(list []rune, text rune) bool {
	for _, v := range list {
		if v == text {
			return true
		}
	}
	return false
}
