package main

import (
    "fmt"
    "time"
)

func main() {
    db, err := NewConnection("db_admin", "db_admin", "db_admin", "127.0.0.1", 5432);
    if err != nil {
        fmt.Println("Connection error");
        panic(err)
    }

    err = AddUser(db, "smss.lite@gmail.com", "+9891234567891", true, "Mahdi", "Shobeiri", "SAG");
    if err != nil {
        fmt.Println("Failed to add user");
        panic(err)
    }
    
    user := GetUser(db, "smss.lite@gmail2fs.com")
    fmt.Printf("User: %v\n", user); /// User: <nil>

    user = GetUser(db, "smss.lite@gmail.com")
    fmt.Printf("User: %v\n", user); /// User: &{6 smss.lite@gmail.com +9891234567891 true Mahdi Shobeiri SAG}

    token := GetUnauthorizedToken(db, "sabzi")
    fmt.Printf("Token: %v\n", token); /// Token: <nil>

    exp := time.Now().Add(time.Duration(20) * time.Minute)
    err = AddUnauthorizedToken(db, user.ID, "sabzi", exp)
    if err != nil {
        fmt.Println("Failed to add token");
        panic(err)
    }

    token = GetUnauthorizedToken(db, "sabzi")
    fmt.Printf("Token: %v\n", token); /// Token: &{6 {0   false   } sabzi 2022-10-17 21:26:21.514783 +0330 +0330}
}