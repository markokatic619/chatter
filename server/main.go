package main

import (
	"bytes"
	databasehandler "chatter/databaseHandler"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type LoginCookies struct {
	id           int
	loginCookie  string
	clientAdress string
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterForm struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Birthday  string `json:"birthday"`
	Password  string `json:"password"`
}
type LogoutForm struct {
	loginCookie string `json:"loginCookie`
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
			//send message to user if user is loged in
			sendMessage(userId, message.Message, &loginCookies)
			//save message to database

			//add user that sent message to top of receivedMessages row

		}

	})
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		go login(w, r, &loginCookies)

	})
	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		go logout(r, &loginCookies)

	})
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		go register(w, r, &loginCookies)

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

func login(w http.ResponseWriter, r *http.Request, loginCookies *[]LoginCookies) {

	if r.Method == http.MethodPost {
		var user LoginForm

		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &user)
		userId := databasehandler.GetUserId(user.Email, user.Password)
		if userId != -1 {
			var cookie LoginCookies
			cookie.loginCookie = StringWithCharset(32)
			cookie.id = userId
			cookie.clientAdress = r.RemoteAddr
			*loginCookies = append(*loginCookies, cookie)
			w.Write([]byte("{\"loginCookie\":\"" + cookie.loginCookie + "\"}"))
		} else {
			w.Write([]byte("{\"loginCookie\":\"\"}"))
		}
	}
}

func register(w http.ResponseWriter, r *http.Request, loginCookies *[]LoginCookies) {
	if r.Method == http.MethodPost {
		var user RegisterForm

		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &user)

		//check is email being used

		if databasehandler.EmailExists(user.Email) {
			w.Write([]byte("{\"loginCookie\":\"\",\"responseMessage\":\"Error: email already in use\"}"))
			return
		}
		//generate username
		var username = ""
		var i = 0
		for ok := true; ok; ok = databasehandler.UsernameExists(username) {
			username = user.FirstName + user.LastName + fmt.Sprint(i)
			i++
		}
		//register new user
		databasehandler.AddUser(username, user.Password, user.Email, user.FirstName, user.LastName, user.Birthday, "")
		//set cookie for new user and return login cookie and response message of success or failure
		userId := databasehandler.GetUserId(user.Email, user.Password)
		if userId != -1 {
			var cookie LoginCookies
			cookie.loginCookie = StringWithCharset(32)
			cookie.id = userId
			cookie.clientAdress = r.RemoteAddr
			*loginCookies = append(*loginCookies, cookie)
			w.Write([]byte("{\"loginCookie\":\"" + cookie.loginCookie + "\",\"responseMessage\":\"success\"}"))
		} else {
			w.Write([]byte("{\"loginCookie\":\"\",\"responseMessage\":\"Error: failed to register\" }"))
		}
	}
}

func remove[T comparable](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

func logout(r *http.Request, loginCookies *[]LoginCookies) {
	if r.Method == http.MethodPost {

		var logout LogoutForm
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &logout)
		for i := 0; i < len(*loginCookies); i++ {
			if (*loginCookies)[i].loginCookie == logout.loginCookie {
				remove((*loginCookies), i)
				return
			}
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
func sendMessage(userId int, message string, cookies *[]LoginCookies) {
	client := http.Client{}
	buf := bytes.NewBufferString(message)
	clientIpAdresses := getClientAddrsForUserId(userId, cookies)
	for _, adress := range clientIpAdresses {
		req, err := http.NewRequest("POST", adress, buf)
		if err != nil {
			// Handle error
		}
		resp, err := client.Do(req)
		if err != nil {
			// Handle error
		}
		defer resp.Body.Close()
	}
}

func getClientAddrsForUserId(id int, cookies *[]LoginCookies) []string {
	var clientAddrs []string
	for _, cookie := range *cookies {
		if cookie.id == id {
			clientAddrs = append(clientAddrs, cookie.clientAdress)
		}
	}
	return clientAddrs
}
