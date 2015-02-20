package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	accessToken string = ""
	userId      string = ""
	expiresIn   string = "0"
	chat_id     string = "1"
	admin_id    string = ""
)
var api Api = Api{
	AccessToken: accessToken,
	UserId:      userId,
	ExpiresIn:   expiresIn,
}

var cache map[string]string = make(map[string]string)

func main() {
	last_msg := 0
	is_pause := 0
	for {
		m := make(map[string]string)
		m["count"] = "1"
		m["out"] = "0"

		response := api.Request("messages.get", m)
		//fmt.Println(response)

		// msg id
		mid_regexp, _ := regexp.Compile("\"mid\":([0-9]+)")
		mid := mid_regexp.FindStringSubmatch(response)[1]

		// msg owner id
		uid_regexp, _ := regexp.Compile("\"uid\":([0-9]+),\"read")
		uid := uid_regexp.FindStringSubmatch(response)[1]

		// msg text
		body_regexp, _ := regexp.Compile("\"body\":\"(.*)\",\"ch")
		body := "no"
		if len(body_regexp.FindStringSubmatch(response)) != 0 {
			body = body_regexp.FindStringSubmatch(response)[1]
		}
		midd, _ := strconv.Atoi(mid)

		if midd > last_msg && body != "no" {
			args := strings.Split(body, " ")

			// pause func
			if body == "!pause" || body == "!пауза" {
				if uid == admin_id {
					if is_pause == 0 {
						is_pause = 1
						send(chat_id, "Пауза включена.")
					} else {
						is_pause = 0
						send(chat_id, "Пауза выключена.")
					}
				} else {
					send(chat_id, "Вы не мой админ, уйдите.")
				}
			}

			fmt.Println(getUserName(uid, "nom") + ":" + mid + " >> " + body)
			if is_pause == 0 {
				switch args[0] {
				case "!ping", "!пинг":
					send(chat_id, "Pong. Your ID -- "+uid)
				case "!uptime", "!аптайм":
					send(chat_id, lexec("uptime"))
				case "!sh", "!к":
					if uid == admin_id && len(args[1:]) != 0 {
						send(chat_id, lexec(strings.Join(args[1:], " ")))
					}
				}
			}

			last_msg = midd
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

func send(chat_id string, msg string) {
	params := make(map[string]string)
	params["chat_id"] = chat_id
	params["message"] = msg

	api.Request("messages.send", params)
}

func getUserName(id string, ncase string) string {
	if len(cache[id]) == 0 {
		params := make(map[string]string)
		params["user_ids"] = id
		params["name_case"] = ncase // nom gen dat acc ins abl

		user := api.Request("users.get", params)

		firstName_regexp, _ := regexp.Compile("\"first_name\":\"(.*)\",\"last")
		lastName_regexp, _ := regexp.Compile("\"last_name\":\"(.*)\"}")
		firstName := firstName_regexp.FindStringSubmatch(user)
		lastName := lastName_regexp.FindStringSubmatch(user)
		cache[id] = firstName[1] + " " + lastName[1]
	}
	return cache[id]
}

func getRandUser() string {
	params := make(map[string]string)
	params["chat_id"] = chat_id

	chatUsers := api.Request("messages.getChatUsers", params)

	users_regexp, _ := regexp.Compile("([0-9]+)")
	users := users_regexp.FindAllStringSubmatch(chatUsers, 50)
	return users[rand.Intn(len(users))][0]
}

func lexec(cmd string) string {
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	return string(out)
}
