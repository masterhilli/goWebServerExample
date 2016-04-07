package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)


/*

+--------------+---------+------+-----+---------+-------+
| Field        | Type    | Null | Key | Default | Extra |
+--------------+---------+------+-----+---------+-------+
| number       | int(11) | NO   | PRI | NULL    |       |
| squareNumber | int(11) | NO   |     | NULL    |       |
+--------------+---------+------+-----+---------+-------+

*/
func main() {
	//startMySQL()
	StartMongoDB()
}


func startMySQL() {
	db, err := sql.Open("mysql", "masterhilli:test@/gowebtest")
	if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	//defer db.Close()

	//dbname := "gowebtest"
	//tablename := "squarenum"
	res, err := db.Exec("DROP TABLE gowebtest.squarenum")
	if (err != nil) {
		fmt.Println("table drop returned an error: " + err.Error())
	} else {
		affected, _ := res.RowsAffected()
		fmt.Printf("Rows affected: %d\n", affected)
	}

	// create table, every time:
	res, err =db.Exec("CREATE TABLE gowebtest.squarenum (number INT(11) NOT NULL, squareNumber INT(11) NULL DEFAULT NULL, PRIMARY KEY (number));");
	if (err != nil) {
		fmt.Println("table create returned an error: " + err.Error())
	} else {
		affected, _ := res.RowsAffected()
		fmt.Printf("Rows affected: %d\n", affected)
	}


	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO squarenum VALUES( ?, ? )") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT squareNumber FROM squarenum WHERE number = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	// Insert square numbers for 0-24 in the database
	for i := 0; i < 25; i++ {
		_, err = stmtIns.Exec(i, (i * i)) // Insert tuples (i, i^2)
		if err != nil {
			//fmt.Printf("Entry already: %d\n", i)
			//panic(err.Error()) // proper error handling instead of panic in your app
		}
	}

	var squareNum int // we "scan" the result in here

	// Query the square-number of 13
	err = stmtOut.QueryRow(13).Scan(&squareNum) // WHERE number = 13
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Printf("The square number of 13 is: %d\n", squareNum)

	// Query another number.. 1 maybe?
	err = stmtOut.QueryRow(1).Scan(&squareNum) // WHERE number = 1
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Printf("The square number of 1 is: %d\n", squareNum)
}