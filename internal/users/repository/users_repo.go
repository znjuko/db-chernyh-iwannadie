package repository

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"main/internal/models"
	"strconv"
)

type UserRepoRealisation struct {
	dbLauncher *sql.DB
}

func NewUserRepoRealisation(db *sql.DB) UserRepoRealisation {
	return UserRepoRealisation{dbLauncher: db}
}

// выдаёт ошибку, если такой пользователь уже существует
func (UserData UserRepoRealisation) CreateNewUser(userModel models.UserModel) ([]models.UserModel, error) {

	allData := make([]models.UserModel,0)

	row , err := UserData.dbLauncher.Query("SELECT nickname , fullname , email , about FROM users WHERE nickname = $1 OR email = $2", userModel.Nickname, userModel.Email)

	if row != nil {
		for row.Next(){

			if err == nil {
				err = errors.New("such user already exists")
			}

			existinUser := models.UserModel{
				Nickname: "",
				Fullname: "",
				Email:    "",
				About:    "",
			}

			row.Scan(&existinUser.Nickname, &existinUser.Fullname, &existinUser.Email, &existinUser.About)

			allData = append(allData,existinUser)

		}

		row.Close()
	}


	if err != nil {
		return allData , errors.New("such user already exists")
	}

	_, err = UserData.dbLauncher.Exec(" INSERT INTO users (nickname , fullname , email , about) VALUES($1 , $2 , $3 ,$4)", userModel.Nickname, userModel.Fullname, userModel.Email, userModel.About)

	allData = append(allData,userModel)

	return allData, err
}

// выдаёт ошибку, если мы пытаемся обновить одно из полей у юзера, которое уже существует
// sql.NoRows - если не существует такого пользователя
func (UserData UserRepoRealisation) UpdateUserData(userModel models.UserModel) (models.UserModel, error) {

	id := 2
	values := make([]interface{},0)

	querySting := "UPDATE users SET"
	nickQuery := " WHERE nickname = $1 RETURNING u_id, nickname, fullname , email, about"
	reqQuery := ""

	values = append(values , userModel.Nickname)

	if userModel.Email != "" {
		values = append(values, userModel.Email)
		reqQuery += " " + "email = $" + strconv.Itoa(id) + ","
		id++
	}

	if userModel.Fullname != "" {
		values = append(values, userModel.Fullname)
		reqQuery += " " + "fullname = $" + strconv.Itoa(id) +","
		id++
	}

	if userModel.About != "" {
		values = append(values, userModel.About)
		reqQuery += " " + "about = $" + strconv.Itoa(id) +","
		id++
	}

	if len(reqQuery) > 1 {
		reqQuery = reqQuery[:len(reqQuery)-1]
	}


	var row *sql.Row

	if len(values) == 1 {
		row = UserData.dbLauncher.QueryRow("SELECT u_id, nickname, fullname , email, about FROM users WHERE nickname = $1", values[0])
	} else {
		row = UserData.dbLauncher.QueryRow(querySting+ reqQuery + nickQuery, values...)
	}


	userId := 0

	err := row.Scan(&userId, &userModel.Nickname, &userModel.Fullname, &userModel.Email, &userModel.About)

	return userModel, err

}

func (UserData UserRepoRealisation) GetUserData(nickname string) (models.UserModel, error) {

	userData := models.UserModel{
		Nickname: "",
		Fullname: "",
		Email:    "",
		About:    "",
	}

	row := UserData.dbLauncher.QueryRow("SELECT nickname , fullname , email, about FROM users WHERE nickname = $1", nickname)

	err := row.Scan(&userData.Nickname, &userData.Fullname, &userData.Email, &userData.About)

	return userData , err
}

func (UserData UserRepoRealisation) Status() (models.Status) {

	statAnsw := new(models.Status)
	row := UserData.dbLauncher.QueryRow("SELECT users , forums , threads , messages FROM counter")
	row.Scan(&statAnsw.User,&statAnsw.Forum,&statAnsw.Thread,&statAnsw.Post)

	return *statAnsw
}

func (UserData UserRepoRealisation) Clear() {

	UserData.dbLauncher.Exec("UPDATE counter SET users = $1 , forums = $2 , threads = $3 , messages = $4",0,0,0,0)

	UserData.dbLauncher.Exec("DELETE FROM users;")
	UserData.dbLauncher.Exec("DELETE FROM forums;")
	UserData.dbLauncher.Exec("DELETE FROM threads;")
	UserData.dbLauncher.Exec("DELETE FROM messages;")
	UserData.dbLauncher.Exec("DELETE FROM messageTU;")
	UserData.dbLauncher.Exec("DELETE FROM voteThreads;")
	UserData.dbLauncher.Exec("DELETE FROM threadUF;")
}