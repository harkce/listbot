package listbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type List struct {
	GroupID  string    `json:"group_id"`
	Title    string    `json:"title"`
	List     []string  `json:"list,omitempty"`
	Multiple bool      `json:"multiple"`
	Element  []Element `json:"element,omitempty"`
}

type Element struct {
	Title string   `json:"title"`
	List  []string `json:"list"`
}

const newItemHelper = "\nGunakan '/add <item>' untuk menambahkan item ke list"

func Retrieve(ID string) (*List, error) {
	l := List{GroupID: ID, Title: "", List: make([]string, 0), Multiple: false, Element: make([]Element, 0)}

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
		log.Println(string(raw))
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

func removeMark(item string) string {
	if strings.HasPrefix(item, "✓ ") || strings.HasPrefix(item, "✗ ") {
		item = strings.TrimPrefix(item, "✓ ")
		item = strings.TrimPrefix(item, "✗ ")
	}
	return item
}

func (l *List) LoadList(ID string) string {
	if len(l.List) == 0 {
		return "List kosong" + newItemHelper
	}

	if l.Title == "" {
		l.Title = "List tanpa judul"
	}
	listString := fmt.Sprintf("%s", l.Title)
	for i, item := range l.List {
		listString = fmt.Sprintf("%s\n%d. %s", listString, i+1, item)
	}
	return listString
}

func (l *List) SetTitle(title string) string {
	l.Title = title
	_, err := save(l.GroupID, *l)
	if err != nil {
		return "Error ganti judul list"
	}
	return fmt.Sprintf("Judul list diganti jadi '%s'", title)
}

func (l *List) AddItem(item string) string {
	l.List = append(l.List, item)
	_, err := save(l.GroupID, *l)
	if err != nil {
		return "Error menambahkan item ke list"
	}
	title := l.Title
	if title == "" {
		title = "list"
	}
	return fmt.Sprintf("Sukses menambahkan '%s' ke %s", item, title)
}

func (l *List) EditItem(pos int, item string) string {
	if len(l.List) == 0 {
		return "List kosong" + newItemHelper
	}

	if pos > len(l.List) || pos < 1 {
		return fmt.Sprintf("List hanya mempunyai %d item", len(l.List))
	}

	l.List[pos-1] = item
	_, err := save(l.GroupID, *l)
	if err != nil {
		return "Error edit list item"
	}
	return fmt.Sprintf("Sukses edit item %d jadi '%s'", pos, item)
}

func (l *List) UncheckItem(pos int) string {
	if len(l.List) == 0 {
		return "List kosong" + newItemHelper
	}

	if pos > len(l.List) || pos < 1 {
		return fmt.Sprintf("List hanya mempunyai %d item", len(l.List))
	}

	l.List[pos-1] = removeMark(l.List[pos-1])
	_, err := save(l.GroupID, *l)
	if err != nil {
		return "Error uncheck item"
	}
	return fmt.Sprintf("Tanda item %d dihilangkan", pos)
}

func (l *List) CheckItem(pos int) string {
	if len(l.List) == 0 {
		return "List kosong" + newItemHelper
	}

	if pos > len(l.List) || pos < 1 {
		return fmt.Sprintf("List hanya mempunyai %d item", len(l.List))
	}

	l.List[pos-1] = "✓ " + removeMark(l.List[pos-1])
	_, err := save(l.GroupID, *l)
	if err != nil {
		return "Error uncheck item"
	}
	return fmt.Sprintf("Item %d ditandai ✓", pos)
}

func (l *List) CrossItem(pos int) string {
	if len(l.List) == 0 {
		return "List kosong" + newItemHelper
	}

	if pos > len(l.List) || pos < 1 {
		return fmt.Sprintf("List hanya mempunyai %d item", len(l.List))
	}

	l.List[pos-1] = "✗ " + removeMark(l.List[pos-1])
	_, err := save(l.GroupID, *l)
	if err != nil {
		return "Error uncheck item"
	}
	return fmt.Sprintf("Item %d ditandai ✗", pos)
}

func (l *List) DeleteItem(pos int) string {
	if len(l.List) == 0 {
		return "List kosong" + newItemHelper
	}

	if pos > len(l.List) || pos < 1 {
		return fmt.Sprintf("List hanya mempunyai %d item", len(l.List))
	}

	deletedItem := l.List[pos-1]
	l.List = append(l.List[0:pos-1], l.List[pos:len(l.List)]...)
	var err error
	_, err = save(l.GroupID, *l)
	if err != nil {
		return "Error hapus item"
	}
	return fmt.Sprintf("Sukses hapus '%s' dari list", deletedItem)
}

func (l *List) ClearItem() string {
	if len(l.List) == 0 {
		return "List kosong" + newItemHelper
	}

	l.List = make([]string, 0)
	l.Title = ""
	_, err := save(l.GroupID, *l)
	if err != nil {
		return "Error hapus list"
	}
	return "List kosong" + newItemHelper
}
