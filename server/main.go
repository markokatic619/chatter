package main

import (
	databasehandler "chatter/databaseHandler"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"log"
	"math/rand"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type LoginCookies struct {
	id          int
	loginCookie string
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Message struct {
	RecieverUsername string `json:"username"`
	Message          string `json:"message"`
}

func main() {

	loginCookies := []LoginCookies{}

	router := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.Use(handlers.CORS(headers, methods, origins))

	router.HandleFunc("/sendMessage", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("loginCookie")
		if err != nil {
			if err == http.ErrNoCookie {
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		userId := getUserIdByCookie(cookie.Value, loginCookies)
		var message Message

		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &message)

		//If user that is sending message is on friend list of a user that is receiving message, send messaage.
		if findFriend(databasehandler.GetUserFriendsList(message.RecieverUsername), databasehandler.GetUserUsernameById(userId)) {

		}

	})
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		setLoginCookie(w, r, &loginCookies)

	})
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("register\n")

	})
	router.HandleFunc("/addFriend", func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("addFriend\n")

	})
	router.HandleFunc("/acceptFriendRequest", func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("acceptFriendRequest\n")

	})

	log.Fatal(http.ListenAndServe(":8080", router))
}

func getUserIdByCookie(cookie string, loginCookies []LoginCookies) int {
	arrLen := len(loginCookies)
	for i := 0; i < arrLen; i++ {
		if cookie == loginCookies[i].loginCookie {
			return loginCookies[i].id
		}
	}
	return -1
}
func setLoginCookie(w http.ResponseWriter, r *http.Request, loginCookies *[]LoginCookies) {

	if r.Method == http.MethodPost {
		var user User

		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &user)
		fmt.Print("Username: " + user.Email + ", password: " + user.Password + "\n")
		userId := databasehandler.GetUserId(user.Email, user.Password)
		if userId != -1 {
			var cookie LoginCookies
			cookie.loginCookie = StringWithCharset(32)
			cookie.id = userId
			*loginCookies = append(*loginCookies, cookie)
			w.Write([]byte("{\"loginCookie\":\"" + cookie.loginCookie + "\"}"))
		}
	}
}
func StringWithCharset(length int) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func findFriend(friendsList string, name string) bool {

	for _, currentName := range strings.Split(friendsList, ",") {
		if currentName == name {
			return true
		}
	}
	return false
}
