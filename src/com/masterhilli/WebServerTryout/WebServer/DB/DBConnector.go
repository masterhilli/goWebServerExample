package db

import (
	. "../UserHandling"
	mgo "gopkg.in/mgo.v2"
	. "../Logger"
)
const (
	TABLE_NAME_USERS string = "WebServerUsers"
	COL_USER_NAME string = "username"
)
type DBConnector struct {
	session *mgo.Session
	index *mgo.Index
	database 	*mgo.Database
	initialized bool
}

func (this *DBConnector) Initialize(dbName string) {
	var err error
	this.session, err = mgo.Dial("localhost:27017")
	if err != nil {
		LOGGER.Printf("Connection to mongo DB was not created: %s\n", err.Error())
	}

	this.index = &mgo.Index{
		Key: []string{COL_USER_NAME},
		Unique: true,
		DropDups: true,
		Background: true, // See notes.
		Sparse: true,
	}

	this.session.SetMode(mgo.Monotonic, true)
	this.database = this.session.DB(dbName)
	c := this.database.C(TABLE_NAME_USERS)
	colInfo := mgo.CollectionInfo{
		Capped : true,
		MaxBytes : 128000}
	err = c.Create(&colInfo)
	if (err == nil) {
		c.EnsureIndex(*this.index)

	} else {
		LOGGER.Printf("ERROR on creation table: %s\n", err.Error())
	}
	err = c.Insert(&User{Username: "admin", Password : "admin"})
	if (err != nil ) {
		LOGGER.Printf("Error occured adding the standard user: Double?(%t) %s \n", mgo.IsDup(err), err.Error())
	}
	this.initialized = true
}


func (this DBConnector) SelectTable(tableName string) TableConnector {
	if this.initialized {
		c := this.database.C(tableName)
		tableConnection := TableConnector{tableLink : c}
		return tableConnection
	} else {
		panic("DBConnector panic: Not yet initialized!\n")
	}
}

func (this *DBConnector) Close() {
	this.session.Close()
}