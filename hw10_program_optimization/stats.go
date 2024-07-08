package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

type DomainStat map[string]int

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

func (us *User) Reset() {
	us.ID = 0
	us.Name = ""
	us.Username = ""
	us.Email = ""
	us.Phone = ""
	us.Password = ""
	us.Address = ""
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

var userPool = sync.Pool{
	New: func() interface{} { return &User{} },
}

func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)
	counter := 0
	for scanner.Scan() {
		user := userPool.Get().(*User)
		user.Reset()
		if err = jsoniter.Unmarshal(scanner.Bytes(), user); err != nil {
			userPool.Put(user)
			return
		}
		result[counter] = *user
		userPool.Put(user)
		counter += 1
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched := strings.Contains(user.Email, "."+domain)

		if matched {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}
