package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"

	"github.com/solenovex/it/funcs"
	"github.com/solenovex/it/model"
)

func registerRoutes() {
	http.HandleFunc("/", listCompanies)
	http.HandleFunc("/companies/seed", seed)
	http.HandleFunc("/companies", listCompanies)
	http.HandleFunc("/companies/add", addCompany)
	http.HandleFunc("/companies/edit/", editCompany)
	http.HandleFunc("/companies/delete/", deleteCompany)
}

func listCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := model.GetAllCompanies()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		funcMap := template.FuncMap{"add": funcs.Add}
		t := template.New("companies").Funcs(funcMap)
		t, err = t.ParseFiles("./templates/_layout.html", "./templates/company/list.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		t.ExecuteTemplate(w, "layout", companies)
	}
}

func addCompany(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t := template.New("company-add")
		t, err := t.ParseFiles("./templates/_layout.html", "./templates/company/add.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		} else {
			t.ExecuteTemplate(w, "layout", nil)
		}
	case http.MethodPost:
		newCompany := model.Company{}
		newCompany.ID = r.PostFormValue("id")
		newCompany.Name = r.PostFormValue("name")
		newCompany.NickName = r.PostFormValue("nickName")
		err := newCompany.Insert()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		} else {
			http.Redirect(w, r, "/companies", http.StatusSeeOther)
		}
	}
}

func editCompany(w http.ResponseWriter, r *http.Request) {
	idPattern := regexp.MustCompile(`/companies/edit/([a-zA-Z0-9]*$)`)
	matches := idPattern.FindStringSubmatch(r.URL.Path)

	if len(matches) > 0 {
		id := matches[1]

		switch r.Method {
		case http.MethodGet:
			company, err := model.GetCompany(id)
			if err == nil {
				t := template.New("company-edit")
				t, err := t.ParseFiles("./templates/_layout.html", "./templates/company/edit.html")
				if err == nil {
					t.ExecuteTemplate(w, "layout", company)
					return
				}
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		case http.MethodPost:
			company := &model.Company{}
			company.ID = r.PostFormValue("id")
			company.Name = r.PostFormValue("name")
			company.NickName = r.PostFormValue("nickName")
			err := company.Update()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			} else {
				http.Redirect(w, r, "/companies", http.StatusSeeOther)
			}
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func deleteCompany(w http.ResponseWriter, r *http.Request) {
	idPattern := regexp.MustCompile(`/companies/delete/([a-zA-Z0-9]*$)`)
	matches := idPattern.FindStringSubmatch(r.URL.Path)

	if len(matches) > 0 {
		id := matches[1]

		if r.Method == http.MethodDelete {
			err := model.DeleteCompany(id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			http.Redirect(w, r, "/companies", http.StatusSeeOther)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func seed(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Ok")
}
