package db

import (
	. "../UserHandling"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	. "../Logger"
	"strings"
)

type TableConnector struct {
	tableLink *mgo.Collection
}

func (this TableConnector) ReceiveStringWhere(colName, value string) string {
	result := User{}
	colName = strings.ToLower(colName)
	LOGGER.Printf("Column Name we search: %s / value of username: %s\n", colName, value)
	err := this.tableLink.Find(bson.M{colName : value}).One(&result)
	if err != nil {
		LOGGER.Printf("TableConnector Find received an error: %s (%s := %s)\n", err.Error(), colName, value)
	}
	return result.Password
}



