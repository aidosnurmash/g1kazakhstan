package models

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

type Database struct {
	Db *sql.DB
}

type NotFoundError struct {
	class_name string
}
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%+v Not Found", e.class_name)
}

func (database *Database) Init() {
	var err error
	database.Db, err = sql.Open("sqlite3", "./pictures.db")
	if err != nil {
		log.Fatal(err)
	}
	sqlStmt := `create table IF NOT EXISTS picture (id integer not null primary key AUTOINCREMENT, path text, origin text, created_at DATETIME DEFAULT CURRENT_TIMESTAMP);`
	_, err = database.Db.Exec(sqlStmt)
	if err != nil {
		fmt.Printf("%q: %s\n", err, sqlStmt)
		//return
	}
	sqlStmt =`create table IF NOT EXISTS part (id integer not null primary key AUTOINCREMENT, path text, picture_id integer not null, part_num integer not null, created_at DATETIME DEFAULT CURRENT_TIMESTAMP);`
	_, err = database.Db.Exec(sqlStmt)
	if err != nil {
		fmt.Printf("%q: %s\n", err, sqlStmt)
		//return
	}

}

func (database *Database) InsertPicture(path, origin string) (int64, error) {
	stmt, err := database.Db.Prepare("insert into picture(path, origin) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(path, origin)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func (database *Database) InsertPart(path string, pictureId int64, partNum int64) (int64, error) {
	stmt, err := database.Db.Prepare("insert into part(path, picture_id, part_num) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(path, pictureId, partNum)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	return id, err
}

func (database *Database) GetAllPictures() []Picture{
	var pictures []Picture
	rows, err := database.Db.Query("select id, path, origin, created_at from picture")
	defer rows.Close()
	for rows.Next() {
		var id int64
		var path string
		var origin string
		var createdAt time.Time

		err = rows.Scan(&id, &path, &origin, &createdAt)
		if err != nil {
			continue
		}
		pictures = append(pictures, Picture{id, path, origin, createdAt})
	}
	return pictures
}

func (database *Database) GetPictureById(query_id int64) (Picture, error) {
	var picture Picture
	stmt, err := database.Db.Prepare("select id, path, origin, created_at from picture where id = ?")
	if err != nil {
		return picture, &NotFoundError{"Picture"}
	}
	defer stmt.Close()
	var id int64
	var path string
	var origin string
	var createdAt time.Time
	err = stmt.QueryRow(strconv.FormatInt(query_id, 10)).Scan(&id, &path, &origin, &createdAt)
	if err != nil {
		return picture, err
	}
	picture = Picture{id, path, origin, createdAt}
	return picture, nil
}

func (database *Database) GetPartById(queryPictureId int64, queryPartNum int) (Part, error) {
	var part Part
	stmt, err := database.Db.Prepare("select id, path, picture_id, part_num, created_at from part where picture_id = ? and part_num = ?")
	if err != nil {
		fmt.Println(err)
		return part, &NotFoundError{"Picture"}
	}
	defer stmt.Close()
	var id int64
	var path string
	var pictureId int64
	var partNum int
	var createdAt time.Time

	row := stmt.QueryRow(strconv.FormatInt(queryPictureId, 10), strconv.Itoa(queryPartNum))

	err = row.Scan(&id, &path, &pictureId, &partNum, &createdAt)
	fmt.Println(id)
	if err != nil {
		fmt.Println(err)
		return part, err
	}
	part = Part{id, pictureId, partNum, path, createdAt}
	return part, nil
}
func (database *Database) GetAllParts() []Part{
	var parts []Part
	rows, err := database.Db.Query("select id, path, picture_id, part_num, created_at from part")
	defer rows.Close()
	for rows.Next() {
		var id int64
		var path string
		var pictureId int64
		var partNum int
		var createdAt time.Time

		err = rows.Scan(&id, &path, &pictureId, &partNum, &createdAt)
		if err != nil {
			continue
		}
		parts = append(parts, Part{id, pictureId, partNum, path, createdAt})
	}
	return parts
}

