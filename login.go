package main


import (
  "net/http"
  "fmt"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "html/template"
  "encoding/json"
  "github.com/gorilla/sessions"
)
var store = sessions.NewCookieStore([]byte("rocket-science"))
var err error
type User struct { 
    UserName    string   `json:"username"`
    FirstName string 	`json:"firstname"`
    Lastname   string   `json:"lastname"`
	Email   string   `json:"email"`
	Contact   string   `json:"contact"`
	Address   string   `json:"address"`
}
type LoginUser struct { 
    UserName    string   `json:"username"`
    Password   string   `json:"password"`
}
func logout(res http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "session-name")
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        return
    }
	// Revoke users authentication
	
	session.Values["authenticated"] = false
	session.Save(req, res)
	res.Write([]byte("Logged Out sucessfully"))
}
func loginPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "login.html")
		return
	}
	session, err := store.Get(req, "session-name")
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        return
    }
	if flashes := session.Flashes(); len(flashes) > 0 {
        // Use the flash values.
		fmt.Println("flashes ",flashes);
    }
	username := req.FormValue("username")
	password := req.FormValue("password")
	 info := &mgo.DialInfo{
        Addrs:    []string{hosts},
        Database: database,
    }

    s, err1 := mgo.DialWithInfo(info)
    if err1 != nil {
        panic(err1)
    }

    col := s.DB(database).C("login")
	
	user := LoginUser {}
	errdb := col.Find(  bson.M{ "username": username, "password": password } ).One(&user)
	 if errdb != nil {
			//http.Error(res, "Invalid UserName and password", 500)
			session.AddFlash("Invalid UserName and password !")
			session.Save(req, res)
            fmt.Println("User not Found ", errdb)
			http.Redirect(res, req, "/login", 301)
            return
        } else{
		session.Values["authenticated"] = true
		session.Save(req, res)
		http.Redirect(res, req, "/home", 301)
		}
		
		

}
func addUserPage(res http.ResponseWriter, req *http.Request) {
session, err := store.Get(req, "session-name")
if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        return
    } 
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(res, "Forbidden. Please login to continue", http.StatusForbidden)
		return
	}

if req.Method != "POST" {
        http.ServeFile(res, req, "addUser.html")
        return
	}
	username := req.FormValue("username")
	firstname := req.FormValue("firstname")
	lastname := req.FormValue("lastname")
	email := req.FormValue("email")
	contact := req.FormValue("contact")
	address := req.FormValue("address")
	 info := &mgo.DialInfo{
        Addrs:    []string{hosts},
        Database: database,
    }

    s, err1 := mgo.DialWithInfo(info)
    if err1 != nil {
        panic(err1)
    }
	col := s.DB(database).C(collection)
	errdb := col.Insert(&User{username, firstname, lastname, email, contact, address})
	if errdb != nil {
			http.Error(res, "Some error occured . Unable to add User", 500)
            fmt.Println("Failed to add Data ", errdb)
            return
        }  
		fmt.Println("Added User Sucessfully")
		//res.Write([]byte("Added User Sucessfully"))
		http.Redirect(res, req, "/home", 301)
}


const (
	hosts = "localhost:27017"
	database = "invoice_revenue"
	collection = "user"
)
func editUserHandler(res http.ResponseWriter, req *http.Request) {
    //io.WriteString(res, "username: "+req.FormValue("username"))
	session, err := store.Get(req, "session-name")
if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        return
    } 
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(res, "Forbidden. Please login to continue", http.StatusForbidden)
		return
	}
	var username = req.FormValue("username")
	  info := &mgo.DialInfo{
        Addrs:    []string{hosts},
        Database: database,
    }

    s, err1 := mgo.DialWithInfo(info)
    if err1 != nil {
        panic(err1)
    }

    col := s.DB(database).C(collection)
	user := User {}
	if req.Method != "POST" {
 
	//fmt.Println("username",username)
	if username != "" {

	
	
	err := col.Find(  bson.M{ "username": username} ).One(&user)
	 if err != nil {
			http.Error(res, "User not found", 500)
            fmt.Println("Failed to fetch Data ", err)
            return
        }  

		tmpl, err := template.ParseFiles("editUser.html")
		err = tmpl.Execute(res, user )
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			
		}
		//fmt.Println("not post")
        return
	}

		
  }
    name := req.FormValue("username")
	firstname := req.FormValue("firstname")
	lastname := req.FormValue("lastname")
	email := req.FormValue("email")
	contact := req.FormValue("contact")
	address := req.FormValue("address")
	colQuerier := bson.M{"username": name}
	change := bson.M{"$set": bson.M{"firstname":firstname,"lastname":lastname,"email":email,"contact":contact,"address":address}}
	
	err = col.Update(colQuerier, change)
	
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		fmt.Println("Failed to add Data ", err)
	}else{
	fmt.Println("Edit sucessful")
	}
	http.Redirect(res, req, "/home", 301)
  
}

func main() {

	
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/addUser", addUserPage)
	 http.HandleFunc("/editUser", editUserHandler)
	 http.HandleFunc("/logout", logout)

	http.HandleFunc("/home", func (res http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "session-name")
if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        return
    } 
	// Check if user is authenticated
	
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
	
		http.Error(res, "Forbidden. Please login to continue", http.StatusForbidden)
		return
	}
	tmpl, err := template.ParseFiles("home.html")
	if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	 info := &mgo.DialInfo{
        Addrs:    []string{hosts},
        Database: database,
    }
	s, err1 := mgo.DialWithInfo(info)
    if err1 != nil {
        panic(err1)
    }

    col := s.DB(database).C(collection)
	var users []User
	err2 := col.Find( bson.M{} ).All(&users)
	 if err2 != nil {
			http.Error(res, "Data not found", 500)
            fmt.Println("Failed to fetch Data ", err)
            return
        }
		pj, err := json.Marshal(users)
	err = tmpl.Execute(res, string(pj) )
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
    

	
})
	http.ListenAndServe(":8080", nil)
	
}