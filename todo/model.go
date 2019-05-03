package todo

import "time"

type Todo struct {
	ID 			string		`json:"id"`
	UserName	string		`json:"username"`
	Text 		string		`json:"text"`
	Complete	bool		`json:"complete"`
	CreatedOn 	time.Time	`json:"created_on"`
}
