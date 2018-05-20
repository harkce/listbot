package listbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type List struct {
	Title string   `json:"title"`
	List  []string `json:"list"`
}

func retrieveListFromDisk(ID string) (*List, error) {
	var l List
	raw, err := ioutil.ReadFile(fmt.Sprintf("%s%s%s.json",
		os.Getenv("GOPATH"),
		"/src/github.com/harkce/listbot/grouplist/",
		ID))
	if err != nil {
		return &l, err
	}

	if err = json.Unmarshal(raw, &l); err != nil {
		return &l, err
	}
	return &l, nil
}

func saveListToDisk(ID string, l List) error {
	raw, err := json.Marshal(l)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s%s%s.json",
		os.Getenv("GOPATH"),
		"/src/github.com/harkce/listbot/grouplist/",
		ID), raw, 0644)

	if err != nil {
		return err
	}
	return nil
}

func LoadList(ID string) string {
	l, err := retrieveListFromDisk(ID)
	if len(l.List) == 0 || err != nil {
		return "List empty"
	}

	if l.Title == "" {
		l.Title = "Untitled list"
	}
	listString := fmt.Sprintf("%s\n", l.Title)
	for i, item := range l.List {
		listString = fmt.Sprintf("%s%d. %s\n", listString, i+1, item)
	}
	return listString
}

func SetTitle(ID, title string) string {
	l, _ := retrieveListFromDisk(ID)
	l.Title = title
	err := saveListToDisk(ID, *l)
	if err != nil {
		return "Error rename list title"
	}
	return fmt.Sprintf("List title changed to %s", title)
}

func AddItem(ID, item string) string {
	l, _ := retrieveListFromDisk(ID)
	l.List = append(l.List, item)
	err := saveListToDisk(ID, *l)
	if err != nil {
		return "Error adding to list"
	}
	title := l.Title
	if title == "" {
		title = "list"
	}
	return fmt.Sprintf("Success add %s to %s", item, title)
}

func DeleteItem(ID string, pos int) string {
	l, err := retrieveListFromDisk(ID)
	if len(l.List) == 0 || err != nil {
		return "List empty"
	}

	if pos > len(l.List) {
		return fmt.Sprintf("List just have %d item(s)", len(l.List))
	}

	deletedItem := l.List[pos-1]
	l.List = append(l.List[0:pos-1], l.List[pos:len(l.List)]...)
	err = saveListToDisk(ID, *l)
	if err != nil {
		return "Error deleting item from list"
	}
	return fmt.Sprintf("Success remove %s from list", deletedItem)
}

func ClearItem(ID string) string {
	l, err := retrieveListFromDisk(ID)
	if len(l.List) == 0 || err != nil {
		return "List empty"
	}

	l.List = []string{}
	l.Title = ""
	err = saveListToDisk(ID, *l)
	if err != nil {
		return "Error clearing list item"
	}
	return "List empty"
}
