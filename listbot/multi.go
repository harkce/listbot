package listbot

import "fmt"

const (
	newListHelper  = "\nGunakan '/newlist <judullist>' untuk membuat list"
	showListHelper = "\n\nGunakan '/list <nomorlist>' untuk melihat isi list"
)

func (l *List) SetMultiple(option string) string {
	if option == "on" {
		l.Multiple = true
	} else if option == "off" {
		l.Multiple = false
	} else {
		return "Format salah"
	}

	if l.Multiple && len(l.List) != 0 && len(l.Element) == 0 {
		l.Element = make([]Element, 0)
		e := Element{Title: l.Title, List: l.List}
		l.Element = append(l.Element, e)
	}

	_, err := save(l.GroupID, *l)
	if err != nil {
		return "Error mengganti opsi multiple list"
	}

	if l.Multiple {
		return "Multiple list diaktifkan"
	}
	return "Multiple list dinonaktifkan"
}

func (l *List) CreateList(title string) string {
	e := Element{Title: title}
	l.Element = append(l.Element, e)
	_, err := save(l.GroupID, *l)
	if err != nil {
		return "Error membuat list baru"
	}

	if title != "" {
		title = " " + title
	}
	return fmt.Sprintf("Sukses membuat list%s, list berada pada posisi %d\nGunakan '/add %d <item>' untuk menambahkan item ke list", title, len(l.Element), len(l.Element))
}

func (l *List) LoadMultiple() string {
	if len(l.Element) == 0 {
		return "List kosong" + newListHelper
	}

	listString := "List"
	for i, e := range l.Element {
		if e.Title == "" {
			e.Title = "List tanpa judul"
		}
		listString = fmt.Sprintf("%s\n%d. %s", listString, i+1, e.Title)
	}
	return listString + showListHelper
}

func (l *List) LoadElement(pos int) string {
	if len(l.Element) == 0 {
		return "List kosong" + newListHelper
	}

	if pos > len(l.Element) || pos < 1 {
		return fmt.Sprintf("Hanya ada %d list di grup ini", len(l.Element))
	}

	e := l.Element[pos-1]
	if e.Title == "" {
		e.Title = "List tanpa judul"
	}
	listString := fmt.Sprintf("%s", e.Title)
	if len(e.List) == 0 {
		listString = fmt.Sprintf("%s\n%s masih kosong\nGunakan '/add %d <item>' untuk menambahkan item", listString, listString, pos)
	} else {
		for i, item := range e.List {
			listString = fmt.Sprintf("%s\n%d. %s", listString, i+1, item)
		}
	}
	return listString
}

func (l *List) SetElementTitle(pos int, title string) string {
	if len(l.Element) == 0 {
		return "List kosong" + newListHelper
	}

	if pos > len(l.Element) || pos < 1 {
		return fmt.Sprintf("Hanya ada %d list di grup ini", len(l.Element))
	}

	e := &l.Element[pos-1]
	e.Title = title
	if _, err := save(l.GroupID, *l); err != nil {
		return "Error ganti judul list"
	}
	return fmt.Sprintf("Judul list nomor %d diganti jadi '%s'", pos, title)
}

func (l *List) AddElementItem(pos int, item string) string {
	if len(l.Element) == 0 {
		return "List kosong" + newListHelper
	}

	if pos > len(l.Element) || pos < 1 {
		return fmt.Sprintf("Hanya ada %d list di grup ini", len(l.Element))
	}

	e := &l.Element[pos-1]
	e.List = append(e.List, item)
	_, err := save(l.GroupID, *l)
	if err != nil {
		return "Error menambahkan item ke list"
	}
	title := e.Title
	if title == "" {
		title = "list"
	}
	return fmt.Sprintf("Sukses menambahkan '%s' ke %s", item, e.Title)
}

func (l *List) EditElementItem(listpos int, pos int, item string) string {
	if len(l.Element) == 0 {
		return "List kosong" + newListHelper
	}

	if listpos > len(l.Element) || listpos < 1 {
		return fmt.Sprintf("Hanya ada %d list di grup ini", len(l.Element))
	}

	e := &l.Element[listpos-1]
	if pos > len(e.List) || pos < 1 {
		return fmt.Sprintf("Hanya ada %d item di list nomor %d", len(e.List), listpos)
	}
	e.List[pos-1] = item
	_, err := save(l.GroupID, *l)
	if err != nil {
		return "Error edit list item"
	}
	if e.Title != "" {
		e.Title = " " + e.Title
	}
	return fmt.Sprintf("Sukses edit item %d di list%s jadi '%s'", pos, e.Title, item)
}

func (l *List) UncheckElementItem(listpos int, pos int) string {
	if len(l.Element) == 0 {
		return "List kosong" + newListHelper
	}

	if listpos > len(l.Element) || listpos < 1 {
		return fmt.Sprintf("Hanya ada %d list di grup ini", len(l.Element))
	}

	e := &l.Element[listpos-1]
	if pos > len(e.List) || pos < 1 {
		return fmt.Sprintf("Hanya ada %d item di list nomor %d", len(e.List), listpos)
	}
	e.List[pos-1] = removeMark(e.List[pos-1])
	_, err := save(l.GroupID, *l)
	if err != nil {
		return "Error edit list item"
	}
	if e.Title != "" {
		e.Title = " " + e.Title
	}
	return fmt.Sprintf("Tanda item %d di list%s dihilangkan", pos, e.Title)
}

func (l *List) CheckElementItem(listpos int, pos int) string {
	if len(l.Element) == 0 {
		return "List kosong" + newListHelper
	}

	if listpos > len(l.Element) || listpos < 1 {
		return fmt.Sprintf("Hanya ada %d list di grup ini", len(l.Element))
	}

	e := &l.Element[listpos-1]
	if pos > len(e.List) || pos < 1 {
		return fmt.Sprintf("Hanya ada %d item di list nomor %d", len(e.List), listpos)
	}
	e.List[pos-1] = "✓ " + removeMark(e.List[pos-1])
	_, err := save(l.GroupID, *l)
	if err != nil {
		return "Error edit list item"
	}
	if e.Title != "" {
		e.Title = " " + e.Title
	}
	return fmt.Sprintf("Item %d di list%s ditandai ✓", pos, e.Title)
}

func (l *List) CrossElementItem(listpos int, pos int) string {
	if len(l.Element) == 0 {
		return "List kosong" + newListHelper
	}

	if listpos > len(l.Element) || listpos < 1 {
		return fmt.Sprintf("Hanya ada %d list di grup ini", len(l.Element))
	}

	e := &l.Element[listpos-1]
	if pos > len(e.List) || pos < 1 {
		return fmt.Sprintf("Hanya ada %d item di list nomor %d", len(e.List), listpos)
	}
	e.List[pos-1] = "✗ " + removeMark(e.List[pos-1])
	_, err := save(l.GroupID, *l)
	if err != nil {
		return "Error edit list item"
	}
	if e.Title != "" {
		e.Title = " " + e.Title
	}
	return fmt.Sprintf("Item %d di list%s ditandai ✗", pos, e.Title)
}

func (l *List) DeleteElementItem(listpos int, pos int) string {
	if len(l.Element) == 0 {
		return "List kosong" + newListHelper
	}

	if listpos > len(l.Element) || listpos < 1 {
		return fmt.Sprintf("Hanya ada %d list di grup ini", len(l.Element))
	}

	e := &l.Element[listpos-1]
	if pos > len(e.List) || pos < 1 {
		return fmt.Sprintf("Hanya ada %d item di list nomor %d", len(e.List), listpos)
	}

	deletedItem := e.List[pos-1]
	e.List = append(e.List[0:pos-1], e.List[pos:len(e.List)]...)
	if _, err := save(l.GroupID, *l); err != nil {
		return "Error hapus item"
	}
	return fmt.Sprintf("Sukses hapus '%s' dari list nomor %d", deletedItem, listpos)
}

func (l *List) RemoveList(pos int) string {
	if len(l.Element) == 0 {
		return "List kosong" + newListHelper
	}

	if pos > len(l.Element) || pos < 1 {
		return fmt.Sprintf("Hanya ada %d list di grup ini", len(l.Element))
	}

	deletedList := l.Element[pos-1].Title
	l.Element = append(l.Element[0:pos-1], l.Element[pos:len(l.Element)]...)
	_, err := save(l.GroupID, *l)
	if err != nil {
		return "Error hapus list"
	}
	return fmt.Sprintf("Sukses hapus list %s", deletedList)
}
