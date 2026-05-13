package main

import (
	"bufio"
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

	scanner := bufio.NewScanner(file)

	i := -1
	seenBrowsers := make(map[string]bool, 100)
	foundUsers := &bytes.Buffer{}

	json := jsoniter.ConfigFastest

	user := User{}
	for scanner.Scan() {
		i++
		line := scanner.Bytes()

		user.Name = ""
		user.Email = ""
		user.Browsers = nil
		// fmt.Printf("%v %v\n", err, line)
		if err := json.Unmarshal(line, &user); err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		browsers := user.Browsers

		for _, browser := range browsers {
			if ok := strings.Contains(browser, "Android"); ok {
				isAndroid = true
				seenBrowsers[browser] = true
			}
			if ok := strings.Contains(browser, "MSIE"); ok {
				isMSIE = true
				seenBrowsers[browser] = true
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
