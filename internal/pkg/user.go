package pkg


type User struct {
	Id int `db:"id"`
	Login string `db:"login"`
	Password string `db:"password_hash"`

}