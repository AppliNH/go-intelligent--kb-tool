package main

import (
	"encoding/json"
	"fmt"
	"go-kb/utils"
	"os"

	"go-kb/kvdb"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/manifoldco/promptui"

	"github.com/google/uuid"
)

var currDb map[string]string = make(map[string]string)
var currSearchType string

var step int = 1

var selectedDocuments []map[string]string

func searchInDocuments(documents map[string]string, searchKey string) []string {

	var searchResults []string
	selectedDocuments = []map[string]string{}

	for k, v := range documents {
		if strings.Contains(k, searchKey) {

			var data map[string]string
			json.Unmarshal([]byte(v), &data)

			if strings.Contains(data[currSearchType], searchKey) {
				searchResults = append(searchResults, data[currSearchType])
				selectedDocuments = append(selectedDocuments, data)
			}
		}

	}

	return utils.RemoveDuplicate(searchResults)

}

func searchInSelectedDocuments(search string) []string {

	var searchResults []string

	for _, doc := range selectedDocuments {
		if strings.Contains(doc[currSearchType], search) {
			searchResults = append(searchResults, doc[currSearchType])
		}

	}
	return utils.RemoveDuplicate(searchResults)

}

func suggest(d prompt.Document) []prompt.Suggest {

	var res []string
	var promptSuggest []prompt.Suggest
	d.Text = strings.ToLower(d.Text)
	if d.Text != "" {
		if step == 1 {
			res = searchInDocuments(currDb, d.Text)
		} else {
			res = searchInSelectedDocuments(d.Text)
		}

		for _, v := range res {
			promptSuggest = append(promptSuggest, prompt.Suggest{Text: strings.Title(strings.ToLower(v))})
		}
	}
	return prompt.FilterHasPrefix(promptSuggest, d.GetWordBeforeCursor(), true)

}

func runSuggestMode(chosen string, notChosen string) (string, string) {

	t1 := prompt.Input(chosen+" ? : ", suggest)

	// Next suggestions will be based on the previously found documents
	step = 2

	// Search type => Name or Surname
	currSearchType = notChosen

	t2 := prompt.Input(notChosen+" ? : ", suggest)

	return t1, t2

}

func runSelectMode(chosen string, notChosen string) (string, string) {
	prompt1 := promptui.Prompt{
		Label: chosen + " ?",
	}

	answer1, _ := prompt1.Run()

	searchInDocuments(currDb, strings.ToLower(answer1)) // will fill up selectedDocuments

	var answer2 string

	if len(selectedDocuments) > 0 {

		var selectedItemsOfNotChosen []string
		for _, v := range selectedDocuments {
			selectedItemsOfNotChosen = append(selectedItemsOfNotChosen, strings.Title(strings.ToLower(v[notChosen])))
		}

		selectedItemsOfNotChosen = utils.RemoveDuplicate(selectedItemsOfNotChosen)
		selectedItemsOfNotChosen = append(selectedItemsOfNotChosen, "Enter a new "+notChosen+"...")

		prompt2 := promptui.Select{
			Label: "Please pick one of these",
			Items: selectedItemsOfNotChosen,
		}

		_, answer2, _ = prompt2.Run()

		if answer2 == "Enter a new "+notChosen+"..." {
			prompt3 := promptui.Prompt{
				Label: notChosen + " ?",
			}
			answer2, _ = prompt3.Run()
		}

	} else {
		prompt2 := promptui.Prompt{
			Label: notChosen + " ?",
		}
		answer2, _ = prompt2.Run()
	}

	return answer1, answer2

}

func main() {

	option := os.Args[1]

	db, err := kvdb.InitDB()
	if err != nil {
		panic(err)
	}

	currDb, err = kvdb.ReadAll(db)
	if err != nil {
		panic(err)
	}

	// Name or surname ?
	choices := []string{"Name", "Surname"}

	nameOrSurname := promptui.Select{
		Label: "Name or surname ?",
		Items: choices,
	}

	i, choice, err := nameOrSurname.Run()
	if err != nil {
		panic(err)
	}

	// Setting the notChosen => Name or Surname
	i2 := (i - 1) * (-1) // Index of the notChosen item
	notChosen := choices[i2]

	// Search type in DB => Name or Surname
	currSearchType = choice

	var t1 string
	var t2 string

	if option == "select" {
		t1, t2 = runSelectMode(choice, notChosen)

	} else if option == "suggest" {
		t1, t2 = runSuggestMode(choice, notChosen)
	}

	result := make(map[string]string)

	result[choice] = strings.ToLower(t1)
	result[notChosen] = strings.ToLower(t2)

	jsonString, err := json.Marshal(result)

	uuid := uuid.Must(uuid.NewRandom())

	if err := kvdb.WriteData(db, result["Name"]+result["Surname"]+uuid.String(), string(jsonString)); err != nil {
		panic(err)
	}

	fmt.Printf("You choose %q\n", strings.Title(strings.ToLower(result["Name"]))+" "+strings.Title(strings.ToLower(result["Surname"])))
}
