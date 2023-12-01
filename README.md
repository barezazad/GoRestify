# GoRestify Backend Boilerplate


[![GolangCI](https://golangci.com/badges/github.com/BarezAzad/GoRestify.svg)](https://golangci.com/r/github.com/BarezAzad/GoRestify)
[![codebeat badge](https://codebeat.co/badges/f7ed90cf-4793-4b82-acd3-00fecf4e3817)](https://codebeat.co/projects/github-com-BarezAzad-GoRestify-master)
[![Go Reference](https://pkg.go.dev/badge/github.com/BarezAzad/GoRestify.svg)](https://pkg.go.dev/github.com/BarezAzad/GoRestify)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

Welcome to the GoRestify Backend Boilerplate, a professional and organized starting point for your GoLang backend projects.

![GoRestify](./GoRestify.webp)




## Table of Contents

1. [Introduction](#introduction)
2. [Setup](#setup)
   - [Redis and MySQL with Docker](#setup-redis-and-mysql-with-docker)
   - [Create Travis User (MySQL)](#create-travis-user-mysql)
3. [Running the Project](#running-the-project)
4. [Activities List](#activities-list)
5. [Database Error Handling](#database-error-handling)
6. [Decimal Handling](#decimal-handling)
7. [Dictionary for Translation](#dictionary-for-translation)
8. [Excel File Generation](#excel-file-generation)
9. [Middleware](#middleware)
10. [Package `pkg_err`](#package-pkg_err)
11. [Package `pkg_http`](#package-pkg_http)
12. [Package `pkg_log`](#package-pkg_log)
13. [Redis Connection](#redis-connection)
14. [Package `pkg_sql`](#package-pkg_sql)
15. [Email Configuration](#email-configuration)
16. [Binding and Validation Model](#binding-and-validation-model)



## Introduction

GoRestify Backend Boilerplate is a robust foundation for building GoLang backend applications. It comes with various features and packages to streamline your development process.

## Setup

### Setup Redis and MySQL with Docker

To set up Redis and MySQL for your project using Docker, you can run the following commands:

```bash
# Start MySQL container
docker run --rm --name db-mysql -d -v mysql-data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=88888888 -e TZ='Asia/Baghdad' -p 3306:3306 mysql --innodb_lock_wait_timeout=1000 --innodb_buffer_pool_size=2147483648  --max_allowed_packet=1073741824  --max_connections=2000 --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

# Start Redis container
docker run --rm --name db-redis -d -v redis-data:/usr/local/etc/redis -p 6379:6379 redis:latest
```

### Create Travis User (MySQL)

Execute the following SQL commands to create a Travis user in MySQL:

```sql
CREATE USER 'travis'@'%';GRANT ALL PRIVILEGES ON *.* TO 'travis'@'%' WITH GRANT OPTION;
```

### Running the Project

To run the GoRestify Backend Boilerplate, use the following commands:

```bash
# Run the admin module
source cmd/admin/sample.env && reflex -r '\.go' -s -- sh -c 'go run cmd/admin/main.admin.go'

# Run the user module
source cmd/user/sample.env &&  reflex -r '\.go' -s -- sh -c 'go run cmd/user/main.user.go'
```

### Get List of Activities
#### You can retrieve a list of activities and record them in the database using the following code snippet:
```go
if data["list"], err = activity.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	if data["count"], err = activity.Count(params); err != nil {
		resp.Error(err).JSON()
		return
	}
```

### Handling Database Errors
#### Handle database errors by using the db_error package as shown below:
```go
    db_error.Parse(err error, tableName string, action validator.Action)

    // In the repository
    err = db_error.Parse(err, base_term.Cities, validator.Update)
```

### Working with Decimal Values
#### In cases involving monetary values, you can perform calculations with decimal values using the decimal package:
```go
quantity := decimal.NewFromInt(3)
fee, _ := decimal.NewFromString(".035")
taxRate, _ := decimal.NewFromString(".08875")

subtotal := price.Num().Mul(quantity.Num())
preTax := subtotal.Num().Mul(fee.Num().Add(decimal.NewFromFloat(1)))
total := preTax.Num().Mul(taxRate.Num().Add(decimal.NewFromFloat(1)))

fmt.Println("Subtotal:", subtotal)                      // Subtotal: 408.06
fmt.Println("Pre-tax:", preTax)                         // Pre-tax: 422.3421
fmt.Println("Taxes:", total.Num().Sub(preTax.Num()))    // Taxes: 37.482861375
fmt.Println("Total:", total)                            // Total: 459.824961375
fmt.Println("Tax rate:", total.Num().Sub(preTax.Num()).Div(preTax.Num())) // Tax rate: 0.08875
```

### Dictionary for Translations
#### Translate terms using the dictionary as follows base on toml file:
```go
    dictionary.Translate(lang Lang, str string, params ...interface{})

    city.Name = dictionary.Translate(dictionary.Ku,city.Name)
```


### Generating Excel Files
#### You can generate Excel files using the excel package:
```go
	ex := excel.New("region")
	ex.AddSheet("Regions").
		AddSheet("Summary").
		Active("Regions").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "D", 15.3).
		SetColWidth("C", "C", 80).
		SetColWidth("D", "F", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Regions").
		WriteHeader("ID", "Name", "Created At").
		SetSheetFields("ID", "Name", "CreatedAt").
		WriteData(regions).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=regions-"+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())
```

### Middleware
#### Use middleware functions like APILogger, JwtAuthGuard, and BasicAuthGuard to enhance your application:
```go
// APILogger logs requests and responses using pkg_log
r.Use(middleware.APILogger())

// JwtAuthGuard decodes tokens and retrieves public and private information from the Authorization header
func JwtAuthGuard() gin.HandlerFunc

// BasicAuthGuard decodes basic authentication for third-party apps from the X-Authorization header
func BasicAuthGuard() gin.HandlerFunc

```

### pkg_err
#### The pkg_err package provides error handling utilities. You can create and customize errors using the following methods:
```go

// New return an initiate of the PkgErr
func New(errStr string, code string,data ...interface{}) *PkgErr

// Take initiate the
func Take(err error, code string,data ...interface{}) *PkgErr

// Message append a message to the error
func (*PkgErr).Message(message string, params ...interface{}) *PkgErr

// Custom is used when some value like status code and basic data needs to be appended to the error
func (*PkgErr).Custom(custom customError) *PkgErr

// InvalidParam is used when want to pint to a field which caused the error
func (*PkgErr).InvalidParam(field string, reason string, params ...interface{}) *PkgErr

// Build return an initiate of the Struct
func (*PkgErr).Build() error

// Example:
// customize error
err = pkg_err.Take(err, "E1174655").Custom(pkg_err.ValidationFailedErr).Build()

// new error
err = pkg_err.New(pkg_err.SomethingWentWrong, "E1171379").
			Message(pkg_err.SomethingWentWrong).Custom(pkg_err.InternalServerErr).Build()
```


### pkg_http
#### You can make HTTP calls to a server using the following functions and data structures:
```go
   // Do http call for a server
   func Do(request Request) (err error)

// Request model
type Request struct {
	Method         string // required
	EndPoint       string // required
	Language       dictionary.Lang
	Headers        []Header
	FormData       FormData
	Payload        interface{} // payload of request
	ParsedResponse interface{} // // parsed response data based on provided model
}

// Header model
type Header struct {
	Key   string
	Value string
}

// FormData model
type FormData struct {
	FileKey  string
	FilePath string
	Payload  []PayloadFormData
}

// PayloadFormData .
type PayloadFormData struct {
	Key   string
	Value string
}
```

### pkg_log
#### The pkg_log package is used for logging errors. You can use it as follows: 
```go
    // will log and stop server
	pkg_log.Fatal(err, "couldn't connect redis cache")

    // it will check if there is error log it
    pkg_log.CheckError(err, "error in regions list")
```

### Redis connection
```go
redisCon, err := pkg_redis.ConnectRedis("redis://:@127.0.0.1:6379/0", true)
```

### pkg_sql
```go
// ForeignKey create foreignKey query for table
func ForeignKey(table string, refTable string, field string, refField ...string) (query string)

engine.DB.Exec(pkg_sql.ForeignKey(base_model.ShopTable, base_model.RegionTable, "region_id"))


// how filter work in in query params
   "[eq]":        " = ",
	"[ne]":       " != ",
	"[gt]":       " > ",
	"[lt]":       " < ",
	"[gte]":      " >= ",
	"[lte]":      " <= ",
	"[like]":     " LIKE ",
	"[and]":      " AND ",
	"[or]":       " OR ",
	"[date]":     " DATE ",
	"[date_gte]": " DATE_GTE ",
	"[date_lte]": " DATE_LTE ",

// example
filter=name[eq]"start shop"[and]region[eq]"Region-1"
```

### Email configuration
```go
	// email config
	engine.EmailConfig = utils.ConfigEmail{
		Host:     engine.Envs[core.EmailHost],
		Port:     engine.Envs.ToInt(core.EmailPort),
		Username: engine.Envs[core.EmailUsername],
		Password: engine.Envs[core.EmailPassword],
	}
	go s.Engine.EmailConfig.SendEmail(toEmail, ccEmail, "alias name to From", subject, body, "attachment path")
```


### Binding and Validation Model
#### The bind tags in the model define validation rules for fields. Here are some examples:
```go
required : required field
min=9 : min length
max=10 : max length
gte=10.6 : greater than
lte=10 : less than
one_of=user_status : define a enum for user_status,and add enum to MustBeInTypes in internal/core/core_action/actions.go
contain=a : that field should be contain (a)
password: accepted regex ".{8,}", "[a-z]", "[A-Z]", "[0-9]"
username: accepted regex "^[a-zA-Z0-9\\._-]+$"
phone: accepted regex "^(964[0-9]{10})$"
email : validate email
birthday: accepted date (min age: current_year-90, max age: current_year-10)
pin: accepted regex "^[0-9]+$"
if_exist : if you have this command before a bind tags it mean just validate in case that field exist

// Example model
type Example struct {
	ID        uint          `json:"id,omitempty"`
	RegionID  uint          `gorm:"index:region_id_idx" json:"region_id,omitempty" bind:"required"`
	Name      string        `gorm:"type:varchar(100);not null;unique" json:"name,omitempty" bind:"required"`
	Status    pkg_types.Enum `gorm:"type:varchar(25);" json:"status,omitempty" bind:"one_of=shop_status"`
	Phone     string        `gorm:"type:varchar(100);" json:"phone,omitempty" bind:"create:required|update:min=7,max=10"` // per action binding
	Password  string        `gorm:"type:varchar(100);" json:"password,omitempty" bind:"if_exist,min=8"` // if_exist example
	CreatedAt *time.Time    `gorm:"->;type:timestamp;not null;default:current_timestamp;" json:"created_at,omitempty"`
	UpdatedAt *time.Time    `gorm:"<-:update;type:timestamp;not null;default:current_timestamp;" json:"updated_at,omitempty"`
	Region    string        `gorm:"->;migration:-" json:"region,omitempty" table:"base_regions.name as region"` // example to join field
}
```