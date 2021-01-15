package main

import (
	"encoding/json"
	"fmt"
	"go-kb/utils"
	"math"
	"sort"
	"strings"

	"go-kb/kvdb"
	. "go-kb/models"

	"github.com/boltdb/bolt"
	"github.com/c-bata/go-prompt"
	"github.com/manifoldco/promptui"
)

var currentType string

func suggest(d prompt.Document) []prompt.Suggest {

	if d.Text == "" {
		return nil
	}

	var promptSuggest []prompt.Suggest

	d.Text = strings.ToLower(d.Text)
	// log.Println(d.Text)
	res := searchFirstItem(currentType, d.Text)
	// log.Println(res)
	var resStrings []string
	for _, i := range res {
		resStrings = append(resStrings, i.Value)
	}

	resStrings = utils.RemoveDuplicate(resStrings)

	for _, v := range resStrings {
		promptSuggest = append(promptSuggest, prompt.Suggest{Text: v})
	}
	// log.Println("_____")
	// log.Println(promptSuggest)
	return prompt.FilterHasPrefix(promptSuggest, d.GetWordBeforeCursor(), true)

}

type ConfObject struct {
	Name       Item
	Surname    Item
	Assoc      Association
	Confidence float64
}

var dbItems *bolt.DB
var dbAssoc *bolt.DB

var currentDbItems map[string]Item = make(map[string]Item)
var currentDbAssocs map[string]Association = make(map[string]Association)

var currentConfidences []ConfObject

func saveItemInDb(item Item) {
	jsonString, err := json.Marshal(item)
	if err != nil {
		panic(err)
	}

	if err := kvdb.WriteData(dbItems, item.Type+";"+item.Value, string(jsonString)); err != nil {
		panic(err)
	}
}

func saveAssocInDb(assoc Association) {
	jsonString, err := json.Marshal(assoc)
	if err != nil {
		panic(err)
	}

	if err := kvdb.WriteData(dbAssoc, assoc.Name+";"+assoc.Surname, string(jsonString)); err != nil {
		panic(err)
	}
}

func searchFirstItem(ttype string, search string) []Item {

	var results []Item
	// log.Println(currentConfidences)
	for _, v := range currentConfidences {
		switch ttype {
		case "name":
			if strings.Contains(v.Name.Value, search) {
				results = append(results, v.Name)
				// log.Print("RESULTS")
				// log.Println(results)
			}
		case "surname":
			if strings.Contains(v.Surname.Value, search) {
				results = append(results, v.Surname)
			}

		}

	}

	return results
}

func generateAllCurrItemsANDAssoc() {

	DBItems, err := kvdb.ReadAll(dbItems)
	if err != nil {
		panic(err)
	}
	DBAssocs, err := kvdb.ReadAll(dbAssoc)
	if err != nil {
		panic(err)
	}

	for k, v := range DBItems {
		item := JSONTOItem(v)
		currentDbItems[k] = item
	}

	for k, v := range DBAssocs {
		assoc := JSONTOAssociation(v)
		currentDbAssocs[k] = assoc
	}

}

func calculateAllConfidences() {

	for k, a := range currentDbAssocs {

		innerKeys := strings.Split(k, ";")

		name := innerKeys[0]
		surname := innerKeys[1]

		itemName := currentDbItems["name;"+name]
		itemSurname := currentDbItems["surname;"+surname]

		confObject := ConfObject{
			Name:       itemName,
			Surname:    itemSurname,
			Assoc:      a,
			Confidence: math.Min(itemName.Pk, itemSurname.Pk) * a.Pk,
		}
		fmt.Println(confObject)
		currentConfidences = append(currentConfidences, confObject)

	}

	// Sort descendly by Confidence
	sort.Slice(currentConfidences, func(i, j int) bool {
		return currentConfidences[i].Confidence > currentConfidences[j].Confidence
	})

}

func returnAssociatedConfidence(name string, surname string) ConfObject {
	for _, c := range currentConfidences {

		if c.Name.Value == name && c.Surname.Value == surname {
			return c
		}
	}

	return ConfObject{Confidence: 0}
}

func main() {
	// Init DB

	var err error

	dbItems, err = kvdb.InitDB("go-kb-items")
	if err != nil {
		panic(err)
	}

	dbAssoc, err = kvdb.InitDB("go-kb-assoc")
	if err != nil {
		panic(err)
	}

	generateAllCurrItemsANDAssoc()
	calculateAllConfidences()

	fmt.Println(strings.Repeat("_", 25))
	fmt.Println("Confidences")
	fmt.Println(currentConfidences)
	fmt.Println(strings.Repeat("_", 25))

	// Name or surname ?
	choices := []string{"name", "surname"}

	var answers map[string]string = make(map[string]string)
	var pickedItems map[string]Item = make(map[string]Item)

	nameOrSurname := promptui.Select{
		Label: "Name or surname ?",
		Items: choices,
	}

	i, chosen, err := nameOrSurname.Run()
	if err != nil {
		s := fmt.Sprintf("Error while selecting name or surname %v", err)
		panic(s)
	}

	// Setting the notChosen => Name or Surname
	i2 := (i - 1) * (-1) // Index of the notChosen item
	notChosen := choices[i2]

	currentType = chosen

	answers[chosen] = prompt.Input(chosen+" ? : ", suggest)

	if str := kvdb.ReadData(dbItems, chosen+";"+answers[chosen]); str != "" {
		pickedItems[chosen] = JSONTOItem(str)
	}

	currentType = notChosen

	answers[notChosen] = prompt.Input(notChosen+" ? : ", suggest)

	if str := kvdb.ReadData(dbItems, notChosen+";"+answers[notChosen]); str != "" {
		pickedItems[notChosen] = JSONTOItem(str)
	}
	var item1 Item
	var item2 Item
	var assoc Association

	var foundItem1 bool = false
	var foundItem2 bool = false

	c := returnAssociatedConfidence(answers["name"], answers["surname"])

	if _, ok := pickedItems[chosen]; ok {
		foundItem1 = true
		item1 = pickedItems[chosen]
		item1.Pk = item1.Pk + (1-item1.Pk)*c.Confidence

	} else {
		item1 = Item{
			Type:  chosen,
			Value: answers[chosen],
			Pk:    0.3,
		}
	}

	if _, ok := pickedItems[notChosen]; ok {
		foundItem2 = true
		item2 = pickedItems[notChosen]
		item2.Pk = item2.Pk + (1-item2.Pk)*c.Confidence
	} else {
		item2 = Item{
			Type:  notChosen,
			Value: answers[notChosen],
			Pk:    0.3,
		}
	}

	assoc = Association{
		Name:    strings.ToLower(answers["name"]),
		Surname: strings.ToLower(answers["surname"]),
		Pk:      0.3,
	}

	if foundItem1 && foundItem2 {
		assoc.Pk = assoc.Pk + (1-assoc.Pk)*c.Confidence
	}

	// To lower
	item1.Value = strings.ToLower(item1.Value)
	item2.Value = strings.ToLower(item2.Value)

	saveItemInDb(item1)
	saveItemInDb(item2)
	saveAssocInDb(assoc)

	generateAllCurrItemsANDAssoc()

	fmt.Println(strings.Repeat("_", 25))
	fmt.Println("Items")
	fmt.Println(strings.Repeat("_", 25))
	fmt.Println(currentDbItems)
	fmt.Println(strings.Repeat("_", 25))
	fmt.Println("Associations")
	fmt.Println(strings.Repeat("_", 25))
	fmt.Println(currentDbAssocs)

}
