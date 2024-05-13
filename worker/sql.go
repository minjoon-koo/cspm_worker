package main

import "worker/config"

var sql string = "select display_name, id , description, created_date_time, expiration_date_time from azuread_group"

var sql2 string = "select display_name, mail, department, member_of from azuread_user"

func main() {
	//db.Connection()
	//controller.GetAllUserInfo()
	//controller.
	//info()
	dd := config.XmlParsher("getAllUsers")
	println(dd)

}

func info() {
	//HTTP.SendPost("/IAM/temp", tmp)
}
