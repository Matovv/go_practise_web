// models package provide structs for our service
package models

type User struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Gender string `json:"gender"`
	Age int `json:"age"`
}