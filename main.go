package main

import (
	"encoding/json"
	"fmt"
	"go-kb/utils"

	"go-kb/kvdb"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/c-bata/go-prompt"
	"github.com/manifoldco/promptui"

	"github.com/google/uuid"
)

var db *bolt.DB
var currDb map[string]string = make(map[string]string)
var currSearchType string

var mode int = 1

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

	if d.Text != "" {
		if mode == 1 {
			res = searchInDocuments(currDb, d.Text)
		} else {
			res = searchInSelectedDocuments(d.Text)
		}

		for _, v := range res {
			promptSuggest = append(promptSuggest, prompt.Suggest{Text: v})
		}
	}
	return prompt.FilterHasPrefix(promptSuggest, d.GetWordBeforeCursor(), true)

}

func main() {
	db, err := kvdb.InitDB()
	if err != nil {
		panic(err)
	}

	currDb, err = kvdb.ReadAll(db)
	if err != nil {
		panic(err)
	}

	//fmt.Println(currDb)

	choices := []string{"Name", "Surname"}

	nameOrSurname := promptui.Select{
		Label: "Name or surname ?",
		Items: choices,
	}
	i, choice, err := nameOrSurname.Run()

	i2 := (i - 1) * (-1) // Index of the notChosen item
	notChosen := choices[i2]

	currSearchType = choice
	t1 := prompt.Input(choice+" ? : ", suggest)

	mode = 2

	currSearchType = notChosen

	//log.Println(selectedDocuments)

	t2 := prompt.Input(notChosen+" ? : ", suggest)

	result := make(map[string]string)

	result[choice] = t1
	result[notChosen] = t2

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	jsonString, err := json.Marshal(result)

	uuid := uuid.Must(uuid.NewRandom())

	if err := kvdb.WriteData(db, result["Name"]+result["Surname"]+uuid.String(), string(jsonString)); err != nil {
		panic(err)
	}

	fmt.Printf("You choose %q\n", result["Name"]+" "+result["Surname"])
}
