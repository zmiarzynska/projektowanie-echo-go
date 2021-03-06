package main

import (
	"net/http"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
	"database/sql"
	"fmt"
	"sync"
	"strconv"
	"github.com/facebookgo/grace/gracehttp"

    _ "github.com/mattn/go-sqlite3"
)

//singleton
var once sync.Once 

type singleton struct {
	db *sql.DB
}

var instance *singleton

func getInstance() *singleton {
    if instance == nil {
        once.Do(
            func() {
                fmt.Println("Creating singleton instance. ")
				instance = &singleton{}
				
				
            })
    } else {
        fmt.Println("Singleton already exists. ")
    }

    return instance
}


// item in db
type (
	item struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

var (
	items = map[int]*item{}
	seq   = 1
)


func (s *singleton)open() *sql.DB {
	//db := initDB("myitems.db")
	//migrate(db)
	db, err := sql.Open("sqlite3", "./myitems.db")
	if err != nil {
		fmt.Println("Error occured")
	} 
	s.db = db
	stmt, _ := db.Prepare("CREATE TABLE IF NOT EXISTS items (id INTEGER PRIMARY KEY AUTOINCREMENT, name text NOT NULL)") 

if err != nil {
	fmt.Println("Error in opening DB")
} 
	stmt.Exec() 

	fmt.Println("Created connection with database. ")
	return db
}

func (s *singleton)close() {
	err := s.db.Close();
	fmt.Println("Connection with database is closed. ")

	if err != nil {
		fmt.Println(err.Error())
	}
}


func (s *singleton) queryI(query string, args ...interface{}) (sql.Result){
	stmt, err := s.db.Prepare(query)
	if err != nil {
		fmt.Println(err.Error())
	}
	result, err2 := stmt.Exec(args...)
	if err2 != nil{
		panic(err2)
	}
	return result
}


func (s *singleton) createItem(c echo.Context) error {

	s.open()
	defer s.close()
	u := &item{
		ID: seq }
	if err := c.Bind(u); err != nil {
		return err
	}
	query := "INSERT INTO items(name) VALUES(?)"
	s.queryI(query,u.Name)
	items[u.ID] = u
	seq++
	fmt.Println("Item created.")

	return c.JSON(http.StatusCreated, u)
}

func (s *singleton) deleteItem(c echo.Context) error {
	s.open()
	defer s.close()
	var id int
	var name string 
	u := &item{
		ID: seq }
	requestedId := c.Param("id")

	if err := c.Bind(u); err != nil {
		return err
	}
	query := "DELETE FROM items Where id = ?"
	
	err := s.db.QueryRow(query, requestedId).Scan(&id, &name)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Item deleted.")
	return c.NoContent(http.StatusNoContent)
}


func (s *singleton) getItem(c echo.Context) error {
	fmt.Println("Get item: ")
	s.open()
	defer s.close()
	var id int
	var name string 
	requestedId := c.Param("id")
	u := &item{
		ID: seq }
		
	if err := c.Bind(u); err != nil {
		return err
	}
	err := s.db.QueryRow("SELECT id, name FROM items WHERE id = ?", requestedId).Scan(&id, &name)
	if err != nil {
		fmt.Println(err)
	}
	result := item{ID: id, Name: name}
	return c.JSON(http.StatusOK, result)

} 

func (s *singleton) updateItem(c echo.Context) error {
	s.open()
	defer s.close()
	var id int
	var name string 
	requestedId := c.Param("id")
	str_to_int, err2:= strconv.Atoi(requestedId)
	if err2 != nil {
		return err2
	}
	u := &item{
		ID: str_to_int}
	
	if err := c.Bind(u); err != nil {
		return err
	}
	err := s.db.QueryRow("UPDATE items SET name=? WHERE id=?", u.Name,requestedId).Scan(&name,&id) // tu gdzies jest prroblem chyba
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Item updated. ")
	
	return c.JSON(http.StatusOK, u)
}


func main() {
	
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to item database!")
	})
	
	s := getInstance()
 	const itemId = "/items/:id"
	

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover()) 

	e.GET(itemId, s.getItem)
	e.PUT(itemId, s.updateItem)
	e.POST("/items", s.createItem)
	e.DELETE(itemId, s.deleteItem)
	e.Server.Addr = ":8000"
	e.Logger.Fatal(gracehttp.Serve(e.Server))

	
}


