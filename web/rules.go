package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/luizbranco/waukeen"
)

func (srv *Server) newRule(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	srv.render(w, nil, "new_rule")
}

func (srv *Server) rules(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rules, err := srv.DB.FindRules("")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
			return
		}

		srv.render(w, rules, "rules")
	case "POST":
		t := r.FormValue("type")
		n, err := strconv.Atoi(t)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		rule := &waukeen.Rule{
			AccountID: r.FormValue("account"),
			Type:      waukeen.RuleType(n),
			Match:     r.FormValue("match"),
			Result:    r.FormValue("result"),
		}

		err = srv.DB.CreateRule(rule)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/rules", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (srv *Server) importRules(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		srv.render(w, nil, "import_rules")
	case "POST":
		file, _, err := r.FormFile("rules")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
			return
		}

		rules, err := srv.RuleImporter.Import(file)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
			return
		}

		for _, r := range rules {
			err := srv.DB.CreateRule(&r)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, err)
				return
			}

		}

		http.Redirect(w, r, "/rules", http.StatusFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}