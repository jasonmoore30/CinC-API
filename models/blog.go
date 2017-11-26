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
	UserNum  int    `json:"user_num"`
	ID       int    `json:"id"`
}

//GetPosts ..
func GetPosts() ([]*Post, error) {
	stmt, err := db.Prepare("SELECT * FROM furmpost")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*Post, 0)
	for rows.Next() {
		post := new(Post)
		err := rows.Scan(&post.Title, &post.Body, &post.PostDate, &post.UserNum, &post.ID)
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
	var post = new(Post)
	for row.Next() {
		err = row.Scan(&post.Title, &post.Body, &post.PostDate, &post.UserNum, &post.ID)
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

	stmt, err := db.Prepare("INSERT INTO furmpost (title, body, postDT, usernum) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	ok := myPost.validate()
	if !ok {
		//return error message about validity
	}
	result, err := stmt.Exec(myPost.Title, myPost.Body, myPost.PostDate, myPost.UserNum)
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
	stmt, err := db.Prepare("UPDATE furmpost SET (title=?, body=? WHERE entryID=?")

	result, err := stmt.Exec(myPost.Title, myPost.Body)
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

func (post *Post) validate() bool {
	return true
}
