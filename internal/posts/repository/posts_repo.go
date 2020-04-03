package repository

import (
	"database/sql"
	"fmt"
	"main/internal/models"
)

type PostRepoRealisation struct {
	dbLauncher *sql.DB
}

func NewPostRepoRealisation(db *sql.DB) PostRepoRealisation {
	return PostRepoRealisation{dbLauncher: db}
}

func (PostRepo PostRepoRealisation) GetPost(id int , flags []string) (models.AllPostData, error) {
	msg := new(models.Message)
	answer := models.AllPostData{}
	row := PostRepo.dbLauncher.QueryRow("SELECT M.m_id,M.date,M.message,M.edit,M.parent,MTU.t_id,U.nickname,F.slug FROM messages M INNER JOIN messageTU MTU ON(MTU.m_id=M.m_id) INNER JOIN users U ON(U.u_id=MTU.u_id) INNER JOIN threadUF TUF ON(TUF.t_id=MTU.t_id) INNER JOIN forums F ON(F.f_id=TUF.f_id) WHERE M.m_id = $1",id)
	err := row.Scan(&msg.Id,&msg.Created,&msg.Message,&msg.IsEdited,&msg.Parent,&msg.Thread,&msg.Author,&msg.Forum)

	if err != nil {
		return answer , err
	}


	answer.Post = msg
	for _ , value := range flags{
		switch value {
		case "user":
			author := new(models.UserModel)
			row = PostRepo.dbLauncher.QueryRow("SELECT nickname , fullname , email, about FROM users WHERE nickname = $1", msg.Author)
			err = row.Scan(&author.Nickname,&author.Fullname,&author.Email,&author.About)

			if err != nil {
				fmt.Println(err , "can't find a user")
			}

			answer.Author = author

		case "forum":
			forum := new(models.Forum)
			fId := 0
			row = PostRepo.dbLauncher.QueryRow("SELECT F.f_id,F.slug,F.title,U.nickname FROM forums F INNER JOIN users U ON(U.u_id=F.u_id) WHERE F.slug= $1",msg.Forum)

			err = row.Scan(&fId, &forum.Slug, &forum.Title, &forum.User)

			if err != nil {
				fmt.Println(err , "can't find a forum")
			}

			row = PostRepo.dbLauncher.QueryRow("SELECT COUNT(DISTINCT t_id) AS thread_counter, COUNT(m_id) as message_counter FROM"+
				" threadUF TUF LEFT JOIN messageTU MTU USING(t_id) WHERE TUF.f_id = $1", fId)

			err = row.Scan(&forum.Threads, &forum.Posts)

			answer.Forum = forum

		case "thread":
			thread := new(models.Thread)
			row = PostRepo.dbLauncher.QueryRow("SELECT T.t_id, T.date, T.message, T.title, T.votes, TUF.slug,U.nickname,F.slug FROM threads T INNER JOIN threadUF TUF ON(TUF.t_id=T.t_id) INNER JOIN users U ON(TUF.u_id=U.u_id) INNER JOIN forums F ON(TUF.f_id=F.f_id) WHERE T.t_id = $1",msg.Thread)
			err = row.Scan(&thread.Id,&thread.Created, &thread.Message, &thread.Title, &thread.Votes, &thread.Slug, &thread.Author,&thread.Forum)

			if err != nil {
				fmt.Println(err , "can't find a thread")
			}

			answer.Thread = thread
		}
	}


	return answer , nil
}

func (PostRepo PostRepoRealisation) UpdatePost(updateData models.Message) (models.Message, error) {

	var err error
	var row *sql.Rows
	if updateData.Message != "" {
		row , err = PostRepo.dbLauncher.Query("UPDATE messages SET edit = CASE WHEN message = $1 THEN FALSE ELSE TRUE END , message = $1  WHERE m_id = $2 RETURNING m_id , date , message , edit, parent", updateData.Message, updateData.Id)
	} else {
		row , err = PostRepo.dbLauncher.Query("SELECT m_id , date , message , edit, parent FROM messages WHERE m_id = $1",updateData.Id)
	}

	if err != nil {
		return updateData , err
	}

	row.Next()


	err = row.Scan(&updateData.Id, &updateData.Created, &updateData.Message, &updateData.IsEdited, &updateData.Parent)
	row.Close()
	if err != nil {
		return updateData , err
	}

	row , err = PostRepo.dbLauncher.Query("SELECT MTU.t_id, U.nickname ,F.slug FROM messageTU MTU INNER JOIN users U ON(U.u_id=MTU.u_id) INNER JOIN threadUF TUF ON(MTU.t_id=TUF.t_id) INNER JOIN forums F ON(TUF.f_id=F.f_id) WHERE MTU.m_id = $1", updateData.Id)
	if err != nil {
		fmt.Println(err)
	}
	if row != nil {
		row.Next()
		err = row.Scan(&updateData.Thread,&updateData.Author,&updateData.Forum)
		row.Close()
	}

	return updateData , nil
}