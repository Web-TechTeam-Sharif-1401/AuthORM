# AuthORM

برای اجرای این ریپو باید  دستورات زیر را وارد کنید:

```commandline
go mod tidy
go run .
```

## db_models.go
این فایل شامل مدل اطلاعات کاربر و توکن‌های منقضی شده‌است.

### UserAccount
```go
type UserAccount struct {
	ID int64 `gorm:"primary_key;auto_increment;not_null"`
	Email  string    `gorm:"unique;not null;default:null;uniqueIndex"`
	PhoneNumber  string    `gorm:"unique;not null;default:null;uniqueIndex"`
	Male    bool   `gorm:"type:bool"`
	FirstName  string
	LastName  string
	PasswordHash  string
}
```

### UnauthorizedToken
```go
type UnauthorizedToken struct {
	UserID int64 `gorm:"foreignKey:;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	User UserAccount
	Token  string `gorm:"index;unique"`
	Expiration  time.Time
}
```

## db_api.go
این فایل شامل تابع‌های لازم برای اضافه کردن و خواندن از دیتابیس است

### NewConnection
با این تابع می‌توانید یک اتصال به دیتابیس بگیرید.  
ورودی‌ها به ترتیب اسم دیتابیس، نام کاربری، رمز، آدرس دیتابیس ،و پورت دیتابیس است.  
```go

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
```

### AddUser
با این تابع می‌توانید یک کاربر جدید به دیتابیس اضافه کنید.  
در صورتی که ایمیل  یا شماره تلفن تکراری باشد این تابع خطا بر می‌گرداند.
```go
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
```

### GetUser
با دادن ایمیل می‌توانید کاربر متناظر با آن ایمیل را بگیرید و در صورتی که کاربری با ایمیل مشخص شده وجود نداشته باشد nil بر می‌گرداند.

```go
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
```

### AddUnauthorizedToken &  GetUnauthorizedToken
این دو تابع هم مشابه دو تابع بالا برای توکن‌های منضی شده است.

```go
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
```

## main.go
این فایل برای نمونه استفاده از دیتابیس نوشته شده است.

### برقراری اتصال به دیتابیس
```go
db, err := NewConnection("db_admin", "db_admin", "db_admin", "127.0.0.1", 5432);
if err != nil {
    fmt.Println("Connection error");
    panic(err)
}
```

### اضافه کردن کاربر
```go
err = AddUser(db, "smss.lite@gmail.com", "+9891234567891", true, "Mahdi", "Shobeiri", "SAG");
if err != nil {
    fmt.Println("Failed to add user");
    panic(err)
}
```

### گرفتن کاربر
```go
user := GetUser(db, "abcdsmss.lite@gmail2fs.com")
fmt.Printf("User: %v\n", user); /// User: <nil>

user = GetUser(db, "smss.lite@gmail.com")
fmt.Printf("User: %v\n", user); /// User: &{6 smss.lite@gmail.com +9891234567891 true Mahdi Shobeiri SAG}
```

### توابع مربوط به توکن
```go
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
```
