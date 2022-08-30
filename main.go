package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/omniskop/vitrum/controls" // Die controls Bibliothek wird importiert, damit sie in vit genutzt werden kann
	"github.com/omniskop/vitrum/gui"
	_ "github.com/omniskop/vitrum/vit/std" // Die std Bibliothek wird importiert, damit sie in vit genutzt werden kann
)

func main() {
	fmt.Println("starting...")

	// Erstellen eines neuen Programms
	app := gui.NewApplication()
	// Der logger wird alle Warnungen ausgeben, die zur Laufzeit im Programm auftreten.
	// Z.B. wenn der JavaScript-Code einen Fehler enthält.
	app.SetLogger(log.New(os.Stdout, "app: ", 0))

	// Erstellung eines neuen Fensters mit Einstiegsdatei
	window, err := app.NewWindow("sources/Calculator.vit")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Eine neue Calculator-Struktur wird für dieses Fenster instanziiert
	NewCalculator(window)

	// Das Programm wird gestartet und die Kontrolle wird an Vitrum übergeben
	err = app.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
