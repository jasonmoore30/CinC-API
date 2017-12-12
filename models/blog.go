package models

import (
	"errors"
	"fmt"
)

//Post ..
type Post struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	PostDate string `json:"post_date"`
	UserNum  int    `json:"usernum"`
	ID       int    `json:"id"`
	// Approved int    `json:"approved"`
}

//GetPosts ..
func GetPosts(admin bool) ([]*Post, error) {

	query := ""
	if admin {
		query = "SELECT * FROM furmpost WHERE adApproval=0"
	} else {
		query = "SELECT * FROM furmpost WHERE adApproval <>0"
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blank int
	posts := make([]*Post, 0)
	for rows.Next() {
		post := new(Post)
		err := rows.Scan(&post.Title, &post.Body, &post.PostDate, &post.UserNum, &post.ID, &blank)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	fmt.Println("Returning posts!")
	return posts, nil
}

//GetPost ..
func GetPost(id string) (*Post, error) {
	row, err := db.Query("SELECT * FROM furmpost WHERE entryID=?", id)
	if err != nil {
		return nil, err
	}
	var blank int
	var post = new(Post)
	for row.Next() {
		err = row.Scan(&post.Title, &post.Body, &post.PostDate, &post.UserNum, &post.ID, &blank)
		if err != nil {
			return nil, err
		}
	}
	if post.ID == 0 || post.Title == "" {
		err := errors.New("post with id: " + id + " not found")
		return nil, err
	}
	return post, nil
}

//AddPost ..
func AddPost(myPost *Post) error {

	stmt, err := db.Prepare("INSERT INTO furmpost (title, body, usernum adApproval) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	ok := myPost.validate()
	if !ok {
		//return error message about validity
	}
	//result, err := stmt.Exec(myPost.Title, myPost.Body, myPost.PostDate, myPost.UserNum)
	//Usernum? Pretty sure we decided on a stateless site, we would need a screenname/email field
	//instead of usernum
	result, err := stmt.Exec(myPost.Title, myPost.Body, myPost.UserNum, 0)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	fmt.Printf("New post created with id: %d", id)
	return nil
}

//DeletePost ..
func DeletePost(id string) error {
	stmt, err := db.Prepare("DELETE FROM furmpost WHERE entryID=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

//UpdatePost ..
func UpdatePost(myPost *Post, id string) error {
	stmt, err := db.Prepare("UPDATE furmpost SET title=?, body=? WHERE entryID=?")

	result, err := stmt.Exec(myPost.Title, myPost.Body, id)
	if err != nil {
		return err
	}
	rowNum, err := result.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println("Posts updated: ", rowNum)
	return nil
}

//ApprovePost sets adApproval to true for the specified post
func ApprovePost(id string) error {
	_, err := db.Query("UPDATE furmpost SET adApproval=? WHERE entryID=?", 1, id)
	if err != nil {
		return err
	}
	return nil
}

//func GetUnapprovedPosts() ([]*Post, error) {}

func (post *Post) validate() bool {
	return true
}
