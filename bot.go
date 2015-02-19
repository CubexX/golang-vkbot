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
	accessToken string = "YOUR_TOKEN"
	userId      string = "YOUR_BOT_ID"
	expiresIn   string = "0"
	chat_id     string = "1"
)
var api Api = Api{
	AccessToken: accessToken,
	UserId:      userId,
	ExpiresIn:   expiresIn,
}

func main() {
	last_msg := 0
	for {
		m := make(map[string]string)
		m["count"] = "1"
		m["out"] = "0"

		response := api.Request("messages.get", m)

		mid_regexp, _ := regexp.Compile("\"mid\":([0-9]+)")
		mid := mid_regexp.FindStringSubmatch(response)[1]

		body_regexp, _ := regexp.Compile("\"body\":\"(.*)\",\"ch")
		body := "no"
		if len(body_regexp.FindStringSubmatch(response)) != 0 {
			body = body_regexp.FindStringSubmatch(response)[1]
		}
		midd, _ := strconv.Atoi(mid)

		if midd > last_msg && body != "no" {
			args := strings.Split(body, " ")

			fmt.Println(mid, body, args)
			switch body {
			case "!ping", "!пинг":
				send(chat_id, "Понг.")
			case "!randu":
				send(chat_id, getUserName(getRandUser(), "nom"))
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
	params := make(map[string]string)
	params["user_ids"] = id
	params["name_case"] = ncase // nom gen dat acc ins abl

	user := api.Request("users.get", params)

	firstName_regexp, _ := regexp.Compile("\"first_name\":\"(.*)\",\"last")
	lastName_regexp, _ := regexp.Compile("\"last_name\":\"(.*)\"}")
	firstName := firstName_regexp.FindStringSubmatch(user)
	lastName := lastName_regexp.FindStringSubmatch(user)
	return firstName[1] + " " + lastName[1]
}

func getRandUser() string {
	params := make(map[string]string)
	params["chat_id"] = chat_id

	chatUsers := api.Request("messages.getChatUsers", params)

	users_regexp, _ := regexp.Compile("([0-9]+)")
	users := users_regexp.FindAllStringSubmatch(chatUsers, 50)
	return users[rand.Intn(len(users))][0]
}

func lexec(cmd string) {
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
}
