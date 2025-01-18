package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/rivo/tview"
)

// Define book struct that will hold the book information
type Book struct {
	Name         string `json:"name"`
	ActualPage   int    `json:"actualpage"`
	LastReadDate string `json:"lastread"`
}

var (
	books    = []Book{}
	bookFile = "books.json"
)

// 1 - Load books from json file - function
func loadBooks() {
	if _, err := os.Stat(bookFile); err == nil {
		{
			data, err := os.ReadFile(bookFile)
			if err != nil {
				log.Fatal("Error reading books file! - ", err)
			}

			// decode JSON and save in variable
			json.Unmarshal(data, &books)
		}
	}
}

// 2 - Save Books
func saveBook() {
	data, err := json.MarshalIndent(books, "", " ")
	if err != nil {
		log.Fatal("Error saving books - ", err)
	}

	os.WriteFile(bookFile, data, 0644)
}

// 3 - Delete book
func deleteBook(index int) {
	if index < 0 || index >= len(books) {
		return
	}

	books = append(books[:index], books[index+1:]...)
	saveBook()
}

func main() {
	app := tview.NewApplication()
	loadBooks()

	bookList := tview.NewTextView().SetDynamicColors(true).SetWordWrap(true)

	bookList.SetBorder(true).SetTitle("Livros em andamento")

	refreshBooks := func() {
		bookList.Clear()

		if len(books) == 0 {
			fmt.Fprintln(bookList, "Nenhum livro em andamento")
		} else {
			for index, item := range books {
				fmt.Fprintf(bookList, "[%d] %s (Página atual: %d, Última leitura: %s)\n", index+1, item.Name, item.ActualPage, item.LastReadDate)
			}
		}
	}

	itemNameInput := tview.NewInputField().SetLabel("Nome do livro: ")
	itemStockInput := tview.NewInputField().SetLabel("Última página: ")

	itemIDInput := tview.NewInputField().SetLabel("ID do livro para excluír: ")

	form := tview.NewForm().
		AddFormItem(itemNameInput).
		AddFormItem(itemStockInput).
		AddFormItem(itemIDInput).
		AddButton("Add Item", func() {
			// Get the text input for name and stock
			nameInput := itemNameInput.GetText()
			actualpageInput := itemStockInput.GetText()

			if nameInput != "" && actualpageInput != "" {

				actualpage, err := strconv.Atoi(actualpageInput)
				if err != nil {
					fmt.Fprintln(bookList, "Digite um valor válido")
					return
				}

				books = append(books, Book{Name: nameInput, ActualPage: actualpage, LastReadDate: time.Now().Local().Format("02 de Jan de 2006")})

				saveBook()

				refreshBooks()

				itemNameInput.SetText("")
				itemStockInput.SetText("")
			}
		}).
		AddButton("Delete Item", func() { // Button to delete an item
			idStr := itemIDInput.GetText()

			if idStr == "" {
				fmt.Fprintln(bookList, "Por favor, insira um Id válido.")
				return
			}

			id, err := strconv.Atoi(idStr)
			if err != nil || id < 1 || id > len(books) {
				fmt.Fprintln(bookList, "Id inválido.")
				return
			}

			deleteBook(id - 1)
			fmt.Fprintf(bookList, "Livro [%d] excluído.\n", id)

			refreshBooks()
			itemIDInput.SetText("")
		}).
		AddButton("Sair", func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle("Meus Livros").SetTitleAlign(tview.AlignLeft)

	flex := tview.NewFlex().
		AddItem(bookList, 0, 1, false).
		AddItem(form, 0, 1, true)

	refreshBooks()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
