package databasehandler

import "strconv"

type inboxTable struct {
	messageId int
	userId    int
	message   string
}

func StartNewChat(userIdOne int, userIdTwo int) {
	dbInbox.Query("create table ?(messageId int auto_increment, userId int, message text, primary key (messageId)")
}

func SendMessage(senderId int, recieverId int, message string) {
	dbInbox.Query("insert into ?(userId,message) values(?,?)", getChatName(senderId, recieverId), senderId, message)
}

func GetChat(from int, to int, userIdOne int, userIdTwo int) []inboxTable {
	lastMessage := GetLastmessageId(userIdOne, userIdTwo)
	var inboxArray []inboxTable
	var inb inboxTable
	result, err := dbInbox.Query("select * from ? where messageId > ? and messageId < ?", getChatName(userIdOne, userIdTwo), lastMessage-from, lastMessage-to)
	if err == nil && result.Next() {
		for result.Next() {
			result.Scan(&inb.messageId, &inb.userId, &inb.message)
			inboxArray = append(inboxArray, inb)
		}
		return inboxArray
	} else {
		return nil
	}
}

func GetLastmessageId(userIdOne int, userIdTwo int) int {
	result, err := dbInbox.Query("select messageId from ? order by messageId desc limit 1", getChatName(userIdOne, userIdTwo))
	var res int
	if err == nil && result.Next() {
		result.Scan(&res)
		return res
	} else {
		return 0
	}
}

func getChatName(userIdOne int, userIdTwo int) string {
	var chatTableName string
	idOne := strconv.Itoa(userIdOne)
	idTwo := strconv.Itoa(userIdTwo)
	if userIdOne > userIdTwo {
		chatTableName = idOne + "|" + idTwo
	} else {
		chatTableName = idTwo + "|" + idOne
	}
	return chatTableName
}
