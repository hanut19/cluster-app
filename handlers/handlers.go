package handlers

import (
	"cluster-app/db"
	"cluster-app/middleware"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// var templates = template.Must(template.ParseGlob("templates/*.html"))

var templates *template.Template

func init() {
	// Try both possible template paths
	paths := []string{
		filepath.Join("templates", "*.html"),       // works in runtime
		filepath.Join("..", "templates", "*.html"), // works in tests from /handlers/
	}

	var err error
	for _, path := range paths {
		templates, err = template.ParseGlob(path)
		if err == nil {
			return
		}
	}
	log.Fatalf("Failed to load templates from known paths: %v", err)
}
func Login(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetUserRole(r)
	if ok {
		http.Redirect(w, r, "/portal", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		if len(strings.TrimSpace(username)) == 0 || len(strings.TrimSpace(password)) == 0 {
			data := map[string]interface{}{
				"error": "username and password is required",
			}
			templates.ExecuteTemplate(w, "login.html", data)
			return
		}
		var hash, role string
		err := db.DB.QueryRow("SELECT password, role FROM users WHERE username=$1", username).Scan(&hash, &role)
		if err != nil {
			data := map[string]interface{}{
				"error": "Invaild username and password",
			}
			templates.ExecuteTemplate(w, "login.html", data)
			return
		}
		if hash != password { //bcrypt.CompareHashAndPassword([]byte(string(hash)), []byte(string(password))) != nil {
			data := map[string]interface{}{
				"error": "Invaild username and password",
			}
			templates.ExecuteTemplate(w, "login.html", data)
			return
		}

		http.SetCookie(w, &http.Cookie{Name: "session", Value: role, Path: "/"})
		http.Redirect(w, r, "/portal", http.StatusSeeOther)
		return
	}
	templates.ExecuteTemplate(w, "login.html", nil)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "session", Value: "", Path: "/", MaxAge: -1})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Portal(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.DB.Query("SELECT name, nodes FROM clusters")
	defer rows.Close()

	var clusters []struct {
		Name  string
		Nodes int
	}
	for rows.Next() {
		var name string
		var nodes int
		rows.Scan(&name, &nodes)
		clusters = append(clusters, struct {
			Name  string
			Nodes int
		}{name, nodes})
	}

	role, _ := r.Cookie("session")
	templates.ExecuteTemplate(w, "portal.html", map[string]interface{}{
		"Clusters": clusters,
		"Role":     role.Value,
	})
}
func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()

		// Extract all cluster names
		names := r.Form["names"]

		for _, name := range names {
			// The input name is counts[<clusterName>], so extract the value
			node := r.FormValue("nodes[" + name + "]")
			if node == "" {
				continue // skip if missing
			}

			_, err := db.DB.Exec("UPDATE clusters SET nodes=$1 WHERE name=$2", node, name)
			if err != nil {
				http.Error(w, "Failed to update cluster: "+name, http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(w, r, "/portal", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

/*func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}*/
