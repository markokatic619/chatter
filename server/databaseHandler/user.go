package databasehandler

import (
	"database/sql"
)

func selectOneFromUser(rowName string, inputName string, input interface{}) *sql.Rows {
	result, err := db.Query("select ? from user where ? = ?", rowName, inputName, input)
	if err == nil {
		return result
	} else {
		return nil
	}
}

func returnResult(result *sql.Rows) interface{} {
	var res any
	if result.Next() {
		result.Scan(&res)
		return res
	} else {
		return -1
	}
}

func GetUserId(email string, password string) int {
	result, err := db.Query("select id from user where email = ? and password = ?", email, password)
	if err == nil {
		return int(returnResult(result).(int64))
	} else {
		return -1
	}
}

func GetUserFullName(username string) (string, string) {
	result, err := db.Query("select firstname,lastname from user where username = ?", username)
	var firstname string
	var lastname string
	if result.Next() && err == nil {
		result.Scan(&firstname, &lastname)
		return firstname, lastname
	} else {
		return "", ""
	}
}

func GetUserBirthday(username string) interface{} {
	return returnResult(selectOneFromUser("birthDay", "username", username))
}

func GetUserFriendsList(username string) string {
	return returnResult(selectOneFromUser("friendsList", "username", username)).(string)
}

func GetUserFriendsListById(userId int) string {
	return returnResult(selectOneFromUser("friendsList", "id", userId)).(string)
}

func GetUserGroupsList(userId int) interface{} {
	return returnResult(selectOneFromUser("groupsList", "id", userId))
}

func GetUserProfilePicture(username string) string {
	return returnResult(selectOneFromUser("profilePicture", "username", username)).(string)
}

func GetUserIdByUsername(username string) int {
	return returnResult(selectOneFromUser("id", "username", username)).(int)
}

func GetUserFriendRequests(userId int) string {
	return (returnResult(selectOneFromUser("friendRequests", "id", userId))).(string)
}

func GetUserUsernameById(userId int) string {
	return (returnResult(selectOneFromUser("username", "id", userId))).(string)
}

func AddUserToFriendsList(userId int, newFriend string) {
	friendsList := GetUserFriendsListById(userId) + "," + newFriend
	db.Query("update user set friendsList = ? where id = ?", friendsList, userId)
}

func AddUser(username string, password string, email string, firstName string, lastName string, birthDay string, profilePicture string) {
	db.Query("insert into user(username,password,email,firstName,lastName,birthDay,profilePicture) values(?,?,?,?,?,?,?)", username, password, email, firstName, lastName, birthDay, profilePicture)
}
