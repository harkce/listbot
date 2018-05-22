package listbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type List struct {
	Title string   `json:"title"`
	List  []string `json:"list"`
}

func retrieve(ID string) (*List, error) {
	var l List

	res, err := http.Get(fmt.Sprintf("%s/get/%s", os.Getenv("KV_HOST"), ID))
	if err != nil {
		log.Println("Error seding request:", err)
		return &l, err
	}

	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		return &l, err
	}

	if err = json.Unmarshal(raw, &l); err != nil {
		log.Println("Error unmarshal:", err)
		return &l, err
	}
	return &l, nil
}

func save(ID string, l List) (*List, error) {
	var list List

	raw, err := json.Marshal(l)
	if err != nil {
		return &list, err
	}

	URL := fmt.Sprintf("%s/store/%s", os.Getenv("KV_HOST"), ID)
	res, err := http.Post(URL, "application/json", bytes.NewBuffer(raw))
	if err != nil {
		return &list, err
	}

	raw, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		return &list, err
	}

	if err = json.Unmarshal(raw, &list); err != nil {
		log.Println("Error unmarshal:", err)
		return &list, err
	}
	return &list, nil
}

func UnsetEnv(ID string) error {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/remove/%s", os.Getenv("KV_HOST"), ID), nil)
	if err != nil {
		log.Println("Error create request:", err)
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return err
	}

	if _, err = ioutil.ReadAll(res.Body); err != nil {
		log.Println("Error reading response:", err)
		return err
	}

	return nil
}

func LoadList(ID string) string {
	l, err := retrieve(ID)
	if len(l.List) == 0 || err != nil {
		return "List empty"
	}

	if l.Title == "" {
		l.Title = "Untitled list"
	}
	listString := fmt.Sprintf("%s", l.Title)
	for i, item := range l.List {
		listString += "\n"
		listString = fmt.Sprintf("%s%d. %s", listString, i+1, item)
	}
	return listString
}

func SetTitle(ID, title string) string {
	l, _ := retrieve(ID)
	l.Title = title
	_, err := save(ID, *l)
	if err != nil {
		return "Error rename list title"
	}
	return fmt.Sprintf("List title changed to %s", title)
}

func AddItem(ID, item string) string {
	l, _ := retrieve(ID)
	l.List = append(l.List, item)
	_, err := save(ID, *l)
	if err != nil {
		return "Error adding to list"
	}
	title := l.Title
	if title == "" {
		title = "list"
	}
	return fmt.Sprintf("Success add %s to %s", item, title)
}

func EditItem(ID string, pos int, item string) string {
	l, err := retrieve(ID)
	if len(l.List) == 0 || err != nil {
		return "List empty"
	}

	if pos > len(l.List) {
		return fmt.Sprintf("List just have %d item(s)", len(l.List))
	}

	l.List[pos-1] = item
	_, err = save(ID, *l)
	if err != nil {
		return "Error edit list item"
	}
	return "Sucess edit item"
}

func DeleteItem(ID string, pos int) string {
	l, err := retrieve(ID)
	if len(l.List) == 0 || err != nil {
		return "List empty"
	}

	if pos > len(l.List) {
		return fmt.Sprintf("List just have %d item(s)", len(l.List))
	}

	deletedItem := l.List[pos-1]
	l.List = append(l.List[0:pos-1], l.List[pos:len(l.List)]...)
	if len(l.List) == 0 {
		err = UnsetEnv(ID)
	} else {
		_, err = save(ID, *l)
	}
	if err != nil {
		return "Error deleting item from list"
	}
	return fmt.Sprintf("Success remove %s from list", deletedItem)
}

func ClearItem(ID string) string {
	l, err := retrieve(ID)
	if len(l.List) == 0 || err != nil {
		return "List empty"
	}

	err = UnsetEnv(ID)
	if err != nil {
		return "Error clearing list item"
	}
	return "List empty"
}
