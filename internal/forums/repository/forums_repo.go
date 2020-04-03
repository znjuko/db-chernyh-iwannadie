package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"main/internal/models"
)

type ForumRepoRealisation struct {
	dbLauncher *sql.DB
}

func NewForumRepoRealisation(db *sql.DB) ForumRepoRealisation {
	return ForumRepoRealisation{dbLauncher: db}
}

func (Forum ForumRepoRealisation) CreateNewForum(forum models.Forum) (models.Forum, error) {

	userId := 0

	row := Forum.dbLauncher.QueryRow("SELECT u_id, nickname FROM users WHERE nickname = $1", forum.User)

	err := row.Scan(&userId, &forum.User)

	if err != nil {
		return forum, err
	}

	_, err = Forum.dbLauncher.Exec("INSERT INTO forums (slug , title, u_id) VALUES($1 , $2 , $3)", forum.Slug, forum.Title, userId)

	if err != nil {

		row := Forum.dbLauncher.QueryRow("SELECT nickname , title , slug FROM forums INNER JOIN users USING(u_id) WHERE slug = $1;", forum.Slug)

		row.Scan(&forum.User, &forum.Title, &forum.Slug)
		// сканировать все данные о форуме, если  err != nil
		return forum, err
	}

	return forum, nil
}

func (Forum ForumRepoRealisation) GetForum(slug string) (models.Forum, error) {

	forumData := new(models.Forum)
	fId := 0
	row := Forum.dbLauncher.QueryRow("SELECT f_id , slug , title, nickname FROM forums INNER JOIN users USING(u_id) WHERE slug = $1", slug)

	err := row.Scan(&fId, &forumData.Slug, &forumData.Title, &forumData.User)

	if err != nil {
		return *forumData, err
	}

	row = Forum.dbLauncher.QueryRow("SELECT COUNT(DISTINCT t_id) AS thread_counter, COUNT(m_id) as message_counter FROM"+
		" threadUF TUF LEFT JOIN messageTU MTU USING(t_id) WHERE TUF.f_id = $1", fId)

	err = row.Scan(&forumData.Threads, &forumData.Posts)

	if err != nil {
		fmt.Println(err, "forum details")
	}

	return *forumData, nil
}

func (Forum ForumRepoRealisation) CreateThread(thread models.Thread) (models.Thread, error) {

	userId := 0
	forumId := 0

	row := Forum.dbLauncher.QueryRow("SELECT u_id , nickname FROM users WHERE nickname = $1", thread.Author)

	err := row.Scan(&userId, &thread.Author)

	if err != nil {
		return thread, err
	}

	row = Forum.dbLauncher.QueryRow("SELECT f_id , slug FROM forums WHERE slug = $1", thread.Forum)

	err = row.Scan(&forumId, &thread.Forum)

	if err != nil {
		return thread, err
	}

	if thread.Slug == "" {
		row = Forum.dbLauncher.QueryRow("INSERT INTO threadUF (u_id , f_id) VALUES($1, $2) RETURNING t_id", userId, forumId)
		err = row.Scan(&thread.Id)
	} else {
		row = Forum.dbLauncher.QueryRow("INSERT INTO threadUF (slug , u_id , f_id) VALUES($1, $2, $3) RETURNING t_id, slug", thread.Slug, userId, forumId)
		err = row.Scan(&thread.Id, &thread.Slug)
	}

	if err != nil {
		row = Forum.dbLauncher.QueryRow("SELECT U.nickname,T.date,F.Slug,T.t_id,T.message,TUF.slug,T.title,T.votes FROM threadUF TUF INNER JOIN threads T ON(TUF.t_id=T.t_id) INNER JOIN users U ON(U.u_id=TUF.u_id) INNER JOIN forums F ON(F.f_id=TUF.f_id) WHERE TUF.slug = $1", thread.Slug)
		err = row.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
		return thread, errors.New("thread already exist")
	}
	row = Forum.dbLauncher.QueryRow("INSERT INTO threads (t_id, date , message, title) VALUES($1,$2,$3,$4) RETURNING date", thread.Id, thread.Created, thread.Message, thread.Title)
	row.Scan(&thread.Created)
	return thread, nil
}

func (Forum ForumRepoRealisation) GetThreads(forum models.Forum, limit int, since string, sort bool) ([]models.Thread, error) {

	row := Forum.dbLauncher.QueryRow("SELECT f_id , slug FROM forums WHERE slug = $1", forum.Slug)

	f_id := 0

	err := row.Scan(&f_id, &forum.Slug)

	if err != nil {
		return nil, err
	}

	orderStatus := ""

	var rowThreads *sql.Rows

	if since != "" {

		sorter := ""

		if sort {
			sorter = "<"
			orderStatus = "DESC"
		} else {
			sorter = ">"
			orderStatus = "ASC"
		}
		sinceStatus := "WHERE T.date" + sorter + "=$2" + " "
		rowThreads, err = Forum.dbLauncher.Query("SELECT T.t_id,T.date,T.message,T.title,T.votes,TUF.slug,F.slug,U.nickname FROM threads T INNER JOIN threadUF TUF ON(TUF.t_id=T.t_id) INNER JOIN forums F ON(TUF.f_id=F.f_id) INNER JOIN users U ON(TUF.u_id=U.u_id) "+sinceStatus+"AND F.slug = $3 "+"ORDER BY T.date "+orderStatus+" LIMIT $1", limit, since, forum.Slug)
	} else {

		if sort {
			orderStatus = "DESC"
		} else {
			orderStatus = "ASC"
		}

		rowThreads, err = Forum.dbLauncher.Query("SELECT T.t_id,T.date,T.message,T.title,T.votes,TUF.slug,F.slug,U.nickname FROM threads T INNER JOIN threadUF TUF ON(TUF.t_id=T.t_id) INNER JOIN forums F ON(TUF.f_id=F.f_id) INNER JOIN users U ON(TUF.u_id=U.u_id) "+"WHERE F.slug = $2 "+"ORDER BY T.date "+orderStatus+" LIMIT $1", limit, forum.Slug)
	}

	// select date from threads WHERE date<='2020-12-22 15:33:59.613+03' ORDER BY date DESC LIMIT 10;
	threads := make([]models.Thread, 0)

	if rowThreads != nil {

		for rowThreads.Next() {
			thread := new(models.Thread)

			rowThreads.Scan(&thread.Id, &thread.Created, &thread.Message, &thread.Title, &thread.Votes, &thread.Slug, &thread.Forum, &thread.Author)

			threads = append(threads, *thread)
		}

		rowThreads.Close()
	}

	return threads, nil

}
func (Forum ForumRepoRealisation) GetForumUsers(slug string, limit int, since string, desc bool) ([]models.UserModel ,error) {

	//  переписать на f_id

	checkRow := Forum.dbLauncher.QueryRow("SELECT f_id FROM forums WHERE slug = $1",slug)
	fId :=0
	err := checkRow.Scan(&fId)

	if err != nil {
		return nil , err
	}


	users := make([]models.UserModel,0)
	var row *sql.Rows

	order := "DESC"
	ranger := "<"

	if !desc {
		order = "ASC"
		ranger = ">"
	}

	if since != "" {
		if limit == 0 {
			row , err = Forum.dbLauncher.Query("SELECT DISTINCT U.nickname,U.fullname,U.email,U.about FROM users U LEFT JOIN messageTU MTU ON(MTU.u_id=U.u_id) LEFT JOIN threadUF TUF ON(TUF.u_id=U.u_id OR MTU.t_id=TUF.t_id) WHERE TUF.f_id = $1 AND U.nickname "+ranger+" $2 ORDER BY U.nickname "+ order,fId,since)
		} else {
			row , err = Forum.dbLauncher.Query("SELECT DISTINCT U.nickname,U.fullname,U.email,U.about FROM users U LEFT JOIN messageTU MTU ON(MTU.u_id=U.u_id) LEFT JOIN threadUF TUF ON(TUF.u_id=U.u_id OR MTU.t_id=TUF.t_id) WHERE TUF.f_id = $1 AND U.nickname "+ranger+" $3 ORDER BY U.nickname "+ order+ " LIMIT $2",fId,limit,since)
		}
	} else {
		if limit == 0 {
			row , err = Forum.dbLauncher.Query("SELECT DISTINCT U.nickname,U.fullname,U.email,U.about FROM users U LEFT JOIN messageTU MTU ON(MTU.u_id=U.u_id) LEFT JOIN threadUF TUF ON(TUF.u_id=U.u_id OR MTU.t_id=TUF.t_id) WHERE TUF.f_id = $1 ORDER BY U.nickname "+ order,fId)
		} else {
			row , err = Forum.dbLauncher.Query("SELECT DISTINCT U.nickname,U.fullname,U.email,U.about FROM users U LEFT JOIN messageTU MTU ON(MTU.u_id=U.u_id) LEFT JOIN threadUF TUF ON(TUF.u_id=U.u_id OR MTU.t_id=TUF.t_id) WHERE TUF.f_id = $1 ORDER BY U.nickname "+ order +" LIMIT $2",fId,limit)
		}
	}

	if err != nil {
		fmt.Println(err)
	}

	if row != nil {
		for row.Next() {

			user := new(models.UserModel)
			err = row.Scan(&user.Nickname, &user.Fullname, &user.Email, &user.About)
			users = append(users,*user)
		}

		row.Close()
	}

	return users , nil

}
