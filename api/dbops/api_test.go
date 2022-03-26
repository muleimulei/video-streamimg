package dbops

import (
	"fmt"
	"testing"
	"time"
)

func clearTables() {
	db.Exec("truncate users")
	db.Exec("truncate videoinfo")
	db.Exec("truncate comments")
	db.Exec("truncate sessions")
}

func TestMain(m *testing.M) {
	clearTables()
	m.Run()
	clearTables()
}

func testUserWorkFlow(t *testing.T) {
	t.Run("Add", testAddUser)
	t.Run("Get", testGetUser)
	t.Run("Del", testDelUser)
	t.Run("Reget", testReGetUser)
}

func testAddUser(t *testing.T) {
	err := AddUserCredential("mulei", "123321")
	if err != nil {
		t.Errorf("Error of adduser %v ", err)
	}
}

func testGetUser(t *testing.T) {
	str, err := GetUserCredential("mulei")
	if str != "123321" || err != nil {
		t.Errorf("Error of adduser %v ", err)
		return
	}
}

func testDelUser(t *testing.T) {
	err := DeleteUser("mulei", "123321")
	if err != nil {
		t.Errorf("Error of adduser %v ", err)
	}
}

func testReGetUser(t *testing.T) {
	pwd, err := GetUserCredential("mulei")

	if err == nil {
		t.Errorf("Error of adduser %v ", err)
	}
	if pwd != "" {
		t.Errorf("delete error %v ", err)
	}
}

// comments
func TestComments(t *testing.T) {
	t.Run("Adduser", testAddUser)
	t.Run("AddComments", testAddComments)
	t.Run("ListComments", testListComments)
}

func testAddComments(t *testing.T) {
	vid := "123"
	aid := 1
	content := "I like this video"

	err := AddNewComments(vid, aid, content)
	if err != nil {
		t.Errorf("Error of AddNewComments : %v", err)
	}
}

func testListComments(t *testing.T) {
	vid := "123"

	res, err := ListComments(vid, 0, time.Now().Unix())

	if err != nil {
		t.Errorf("Error of ListComments: %v ", err)
	}

	for i, ele := range res {
		fmt.Printf("Comment %d: %v \n", i, *ele)
	}
}
