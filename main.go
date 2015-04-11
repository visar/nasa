package main

import (
    "github.com/gin-gonic/gin"
    "database/sql"
    "github.com/coopernurse/gorp"
    _ "github.com/mattn/go-sqlite3"
    "log"
    "strconv"
    "time"
)

var dbmap = initDb()

func main(){

    defer dbmap.Db.Close()

    router := gin.Default()
    router.GET("/attributes", AttributesList)
    router.GET("/measures", MeasuresList)
    router.POST("/attributes", AttributePost)
    router.POST("/measures", MeasurePost)
    router.GET("/attributes/:id", AttributesDetail)
    router.GET("/measures/:id", MeasuresDetail)
    router.Run(":8000")
}

type Attribute struct {
    Id int64 `db:"attribute_id"`
    Code string
    Description string
    Unit string
}

type Measure struct {
	MId int64 `db:"measure_id"`
	UserId int64 
	Time int64
	Quantity int64 `json:",string"`
}

func createAttribute(code, description, unit string) Attribute {
    attribute := Attribute{
	Code:	code,
	Description:	description,
	Unit:	unit,
    }

    err := dbmap.Insert(&attribute)
    checkErr(err, "Insert failed")
    return attribute
}

func createMeasure(quantity, userid int64) Measure {
	measure := Measure{
		Time:	time.Now().UnixNano(),
		UserId:	userid,
		Quantity: quantity,
	}

    err := dbmap.Insert(&measure)
    checkErr(err, "Insert failed")
    return measure 
}

func getAttribute(attribute_id int) Attribute {
    attribute := Attribute{}
    err := dbmap.SelectOne(&attribute, "select * from attributes where attribute_id=?", attribute_id)
    checkErr(err, "SelectOne failed")
    return attribute
}

func getMeasure(measure_id int) Measure {
	measure := Measure{}
	err := dbmap.SelectOne(&measure, "select * from measures where measure_id=?", measure_id)
	checkErr(err, "SelectOne failed")
	return measure
}

func AttributesList(c *gin.Context) {
    var attributes []Attribute
    _, err := dbmap.Select(&attributes, "select * from attributes order by attribute_id")
    checkErr(err, "Select failed")
    content := gin.H{}
    for k, v := range attributes {
        content[strconv.Itoa(k)] = v
    }
    c.JSON(200, content)
}

func MeasuresList(c *gin.Context) {
	var measures []Measure
        _, err := dbmap.Select(&measures, "select * from measures order by measure_id")
	checkErr(err, "Select failed")
	content := gin.H{}
	for k, v := range measures {
		content[strconv.Itoa(k)] = v
	}
	c.JSON(200, content)
}

func AttributesDetail(c *gin.Context) {
    attribute_id := c.Params.ByName("id")
    a_id, _ := strconv.Atoi(attribute_id)
    attribute := getAttribute(a_id)
    content := gin.H{"code": attribute.Code, "description": attribute.Description, "unit": attribute.Unit}
    c.JSON(200, content)
}

func MeasuresDetail(c *gin.Context) {
    measure_id := c.Params.ByName("id")
    m_id, _ := strconv.Atoi(measure_id)
    measure := getMeasure(m_id)
    content := gin.H{"userid": measure.UserId, "time": measure.Time, "quantity": measure.Quantity}
    c.JSON(200, content)
}

func AttributePost(c *gin.Context) {
    var json Attribute

    c.Bind(&json) // This will infer what binder to use depending on the content-type header.
    attribute := createAttribute(json.Code, json.Description, json.Unit)
    if attribute.Code == json.Code {
        content := gin.H{
            "result": "Success",
            "code": attribute.Code,
            "description": attribute.Description,
            "unit": attribute.Unit,
        }
        c.JSON(201, content)
    } else {
        c.JSON(500, gin.H{"result": "An error occured"})
    }
}

func MeasurePost(c *gin.Context) {
    var json Measure

    c.Bind(&json) // This will infer what binder to use depending on the content-type header.
    measure := createMeasure(json.Quantity, json.UserId)
    if measure.Quantity == json.Quantity {
        content := gin.H{
            "result": "Success",
            "userid": measure.UserId,
            "qunatity": measure.Quantity,
        }
        c.JSON(201, content)
    } else {
        c.JSON(500, gin.H{"result": "An error occured"})
    }
}

func initDb() *gorp.DbMap {
    db, err := sql.Open("sqlite3", "db.sqlite3")
    checkErr(err, "sql.Open failed")

    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

    dbmap.AddTableWithName(Attribute{}, "attributes").SetKeys(true, "Id")
    dbmap.AddTableWithName(Measure{}, "measures").SetKeys(true, "MId")

    err = dbmap.CreateTablesIfNotExists()
    checkErr(err, "Create tables failed")

    return dbmap
}

func checkErr(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}
