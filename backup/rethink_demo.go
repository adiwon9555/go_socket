package backup

import (
	"fmt"
	"log"

	r "github.com/dancannon/gorethink"
)

type user struct {
	Id   string `gorethink:"id,omitempty"`
	Name string `gorethink:"name"`
}

func main1() {
	session, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "demo",
	})
	if err != nil {
		fmt.Println("Error:- ", err)
		return
	}
	// use := user{
	// 	Name: "Adi",
	// }

	//Insert
	// response, err := r.Table("user").Insert(use).RunWrite(session)
	// if err != nil {
	// 	fmt.Println("Error:- ", err)
	// 	return
	// }
	// fmt.Printf("%#v\n", response)

	//Update
	// response, err := r.Table("user").Get("7efcd1f4-58e4-4cbc-b395-b830219b878c").Update(use).RunWrite(session)
	// if err != nil {
	// 	fmt.Println("Error:- ", err)
	// 	return
	// }
	// fmt.Printf("%#v\n", response)

	//Delete
	// response, err := r.Table("user").Get("30f7dc65-93d3-4a8a-8ceb-be4d9e5bd3b6").Delete().RunWrite(session)
	// if err != nil {
	// 	fmt.Println("Error:- ", err)
	// 	return
	// }
	// fmt.Printf("%#v\n", response)

	//Get one
	// res, err := r.Table("user").Get("23ef58a4-ef66-471e-abbe-406344bf9ea6").Run(session)
	// if err != nil {
	// 	fmt.Println("Error:- ", err)
	// 	return
	// }
	// var us user
	// err = res.One(&us)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Printf("%+v\n", us)

	//Get all
	res, err := r.Table("user").Run(session)
	if err != nil {
		fmt.Println("Error:- ", err)
		return
	}
	var us []user
	err = res.All(&us)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", us)

	//Change feed
	// cursor, err := r.Table("user").Changes(r.ChangesOpts{IncludeInitial: true}).Run(session)
	// if err != nil {
	// 	fmt.Println("Error:- ", err)
	// 	return
	// }
	// var changeResponse r.ChangeResponse
	// for cursor.Next(&changeResponse) {
	// 	fmt.Printf("%#v\n", changeResponse)
	// }
}
