package hw10programoptimization

import (
	"bufio"
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

var userPool = sync.Pool{
	New: func() interface{} { return &User{} },
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	result := make(DomainStat)

	for scanner.Scan() {
		user := userPool.Get().(*User)
		user.Reset()

		if err := jsoniter.Unmarshal(scanner.Bytes(), user); err != nil {
			userPool.Put(user)
			return nil, err
		}

		matched := strings.HasSuffix(user.Email, "."+domain)
		if matched {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}
