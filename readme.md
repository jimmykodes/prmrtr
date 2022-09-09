# Parameterized Router

Parameterized http routing. Doing only what is says on the tin.

---

This is an alternative to the default `net/http` Mux that adds named parameters to routes.

## Usage

```golang
package main

import (
	"net/http"
	
	"github.com/jimmykodes/prmrtr"
)

func main() {
	router := prmrtr.NewRouter()
	router.HandleFunc("/user/:id", func(w http.ResponseWriter, r *http.Request) {
		// will match /user/12 but will also match /user/test since there is currently no way to
		// declare data types per parameter
		vars := prmrtr.Vars(r)
		id, ok := vars.Int("id")
		if !ok {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		// ...
	})
	router.HandleFunc("/user/name/:username", func(w http.ResponseWriter, r *http.Request) {
		// will match /user/name/someUserName
		vars := prmrtr.Vars(r)
		username, ok := vars.String("username")
		if !ok {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		// ...
	})

	// using SubRouters
	todoRouter := router.SubRouter("/todos")
	todoRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// resolves to /todos/
		// ...
	})
	todoRouter.HandleFunc("/:id", func(w http.ResponseWriter, r *http.Request) {
		// resolves to /todos/:id
		// ...
	})
	
	// SubRouters all the way down
	listRouter := todoRouter.SubRouter("/lists")
	itemRouter := listRouter.SubRouter("/items")
	// continue ad nauseam
	
	// the topmost router should be passed to the http server
	http.ListenAndServe(":8080", router)
}
```
