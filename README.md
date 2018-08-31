## gitdb

`gitdb` is a simple git wrapper that lets you use git as a database.

Here's a simple program to get started:

```golang
package main

import (
	"encoding/json"
	"log"

	"github.com/makramkd/gitdb"
)

func main() {
	db := gitdb.GetInstance()
	err := db.Open("./repo1")
	if err != nil {
		log.Fatalf("Could not open repo1: %s", err.Error())
	}

	s := struct {
		Name  string
		Email string
	}{
		Name:  "makram kamaleddine",
		Email: "makram.kd@github.com",
	}
	marshaled, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal struct: %s", err.Error())
	}
	err = db.Save(marshaled, "helloworld.json")
	if err != nil {
		log.Fatalf("Failed to save file to repo: %s", err.Error())
	}
}
```

`gitdb` is very much a dumb proof of concept and not ready for production use.
