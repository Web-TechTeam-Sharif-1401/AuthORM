package main


import (
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
	"time"
)


func NewConnection(name string, user string, pass string, host string, port int) (*gorm.DB, error) {
    dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, pass, host, port, name)
    db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
    if err != nil {
		return nil, err
    }
    db.AutoMigrate(&UserAccount{})
	db.AutoMigrate(&UnauthorizedToken{})
    return db, nil
}


func AddUser(
	db *gorm.DB,
	email  string,
	phone_number  string,
	male    bool ,
	first_name  string,
	last_name  string,
	password_hash  string,
) error{
	return db.Create(&UserAccount{
		Email: email,
		PhoneNumber: phone_number,
		Male: male,
		FirstName: first_name,
		LastName: last_name,
		PasswordHash: password_hash,
	}).Error
}

func GetUser(
	db *gorm.DB,
	email string,
) *UserAccount{
	user := UserAccount{}
    if db.Where("email = ?", email).First(&user).Error != nil {
		return nil
    }
	return &user;
}

func AddUnauthorizedToken(
	db *gorm.DB,
	user_id int64,
	token  string,
	expiration  time.Time,
) error{
	return db.Create(&UnauthorizedToken{
		UserID: user_id,
		Token: token,
		Expiration: expiration,
	}).Error
}

func GetUnauthorizedToken(
	db *gorm.DB,
	token string,
) *UnauthorizedToken{
	result := UnauthorizedToken{}
    if db.Where("token = ?", token).First(&result).Error != nil {
		return nil
    }
	return &result;
}