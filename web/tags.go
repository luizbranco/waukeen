package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/luizbranco/waukeen"
)

func (srv *Server) newTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	srv.render(w, nil, "tag")
}

func (srv *Server) tags(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := r.URL.Path[len("/transactions/"):]
		if id == "" {
			tags, err := srv.DB.FindTags("")

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, err)
				return
			}

			srv.render(w, tags, "tags")
			return
		}

		t, err := srv.DB.FindTag(id)
		if err != nil {
			srv.render(w, nil, "404")
			return
		}
		srv.render(w, t, "tag")
	case "POST":
		b := r.FormValue("budget")
		n, err := strconv.ParseInt(b, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tag := &waukeen.Tag{
			Name:   r.FormValue("name"),
			Budget: n,
		}

		err = srv.DB.CreateTag(tag)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/tags", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}