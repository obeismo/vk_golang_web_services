package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// вам надо написать более быструю оптимальную этой функции

type User struct {
	Name     string
	Email    string
	Browsers []string
}

func FastSearch(out io.Writer) {
	/*
		!!! !!! !!!
		обратите внимание - в задании обязательно нужен отчет
		делать его лучше в самом начале, когда вы видите уже узкие места, но еще не оптимизировалм их
		так же обратите внимание на команду в параметром -http
		перечитайте еще раз задание
		!!! !!! !!!
	*/
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	fileContents, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	seenBrowsers := make([]string, 0)
	uniqueBrowsers := 0
	foundUsers := &bytes.Buffer{}

	lines := bytes.Split(fileContents, []byte("\n"))

	json := jsoniter.ConfigFastest

	users := make([]User, 0, len(lines))
	user := User{}
	for _, line := range lines {
		//if !(bytes.Contains(line, []byte("Android")) && bytes.Contains(line, []byte("MSIE"))) {
		//	continue
		//}
		user.Name = ""
		user.Email = ""
		user.Browsers = nil
		// fmt.Printf("%v %v\n", err, line)
		err := json.Unmarshal(line, &user)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}

	for i, user := range users {

		isAndroid := false
		isMSIE := false

		browsers := user.Browsers

		for _, browser := range browsers {
			if ok := strings.Contains(browser, "Android"); ok {
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		for _, browser := range browsers {
			if ok := strings.Contains(browser, "MSIE"); ok {
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := strings.ReplaceAll(user.Email, "@", " [at] ")
		_, _ = fmt.Fprintf(foundUsers, "[%d] %s <%s>\n", i, user.Name, email)
	}

	_, _ = fmt.Fprintln(out, "found users:")
	_, _ = foundUsers.WriteTo(out)
	_, _ = fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}
