package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"main/internal/models"
	"time"
)

type ThreadRepoRealisation struct {
	dbLauncher *sql.DB
}

func NewThreadRepoRealisation(db *sql.DB) ThreadRepoRealisation {
	return ThreadRepoRealisation{dbLauncher: db}
}

func (Thread ThreadRepoRealisation) CreatePost(slug string, id int, posts []models.Message) ([]models.Message, error) {
	threadId := 0
	var row *sql.Row

	t := time.Now()

	if slug != "" {
		row = Thread.dbLauncher.QueryRow("SELECT t_id FROM threadUF WHERE slug = $1", slug)
	} else {
		row = Thread.dbLauncher.QueryRow("SELECT t_id FROM threadUF WHERE t_id = $1", id)
	}

	err := row.Scan(&threadId)

	if err != nil {
		return nil, err
	}

	currentPosts := make([]models.Message, 0)

	for _, value := range posts {

		if value.Parent == 0 {
			authorId := 0
			value.Thread = threadId

			row = Thread.dbLauncher.QueryRow("SELECT u_id , nickname FROM users WHERE nickname = $1", value.Author)
			err = row.Scan(&authorId, &value.Author)

			if err != nil {
				return []models.Message{value} , errors.New("no user")
			}

			row = Thread.dbLauncher.QueryRow("SELECT F.slug FROM threadUF TUF INNER JOIN forums F ON(TUF.f_id=F.f_id) WHERE TUF.t_id = $1", threadId)
			err = row.Scan(&value.Forum)

			row = Thread.dbLauncher.QueryRow("INSERT INTO messageTU (u_id,t_id) VALUES($1,$2) RETURNING m_id", authorId, threadId)
			row.Scan(&value.Id)
			value.IsEdited = false
			row = Thread.dbLauncher.QueryRow("INSERT INTO messages (m_id, date , message , parent) VALUES ($1 , $2 , $3 , $4) RETURNING date", value.Id, t, value.Message, value.Parent)
			err = row.Scan(&value.Created)
			currentPosts = append(currentPosts, value)
		} else {
			row = Thread.dbLauncher.QueryRow("SELECT M.m_id FROM messages M INNER JOIN messageTU MTU ON(MTU.m_id=M.m_id) WHERE M.m_id = $1 AND MTU.t_id = $2", value.Parent,threadId)

			err = row.Scan(&value.Parent)

			if err != nil {
				return nil, errors.New("Parent post was created in another thread")
			}

			authorId := 0
			value.Thread = threadId

			row = Thread.dbLauncher.QueryRow("SELECT u_id , nickname FROM users WHERE nickname = $1", value.Author)
			err = row.Scan(&authorId, &value.Author)

			if err != nil {
				return []models.Message{value} , errors.New("no user")
			}

			row = Thread.dbLauncher.QueryRow("SELECT F.slug FROM threadUF TUF INNER JOIN forums F ON(TUF.f_id=F.f_id) WHERE TUF.t_id = $1", threadId)
			row.Scan(&value.Forum)

			row = Thread.dbLauncher.QueryRow("INSERT INTO messageTU (u_id,t_id) VALUES($1,$2) RETURNING m_id", authorId, threadId)
			err = row.Scan(&value.Id)
			value.IsEdited = false

			dRow, _ := Thread.dbLauncher.Query("INSERT INTO messages (m_id, date , message , parent) VALUES ($1 , $2 , $3 , $4) RETURNING date", value.Id, t, value.Message, value.Parent)
			dRow.Next()
			dRow.Scan(&value.Created)
			defer dRow.Close()

			currentPosts = append(currentPosts, value)
		}

	}

	return currentPosts, nil
}

func (Thread ThreadRepoRealisation) VoteThread(nickname string, voice, threadId int, thread models.Thread) (models.Thread, error) {

	voterId := 0
	voterNick := ""
	row := Thread.dbLauncher.QueryRow("SELECT u_id , nickname FROM users WHERE nickname = $1", nickname)
	err := row.Scan(&voterId, &voterNick)
	if err != nil {
		return thread, err
	}
	if thread.Slug != "" {
		row = Thread.dbLauncher.QueryRow("SELECT TUF.t_id , TUF.slug, U.nickname , F.slug FROM threadUF TUF INNER JOIN users U USING(u_id) INNER JOIN forums F USING(f_id) WHERE TUF.slug = $1", thread.Slug)
	} else {
		row = Thread.dbLauncher.QueryRow("SELECT TUF.t_id , TUF.slug, U.nickname , F.slug FROM threadUF TUF INNER JOIN users U USING(u_id) INNER JOIN forums F USING(f_id) WHERE TUF.t_id = $1", threadId)
	}

	err = row.Scan(&thread.Id, &thread.Slug, &thread.Author, &thread.Forum)

	if err != nil {
		return thread, err
	}

	voted := 0
	row = Thread.dbLauncher.QueryRow("SELECT counter FROM voteThreads WHERE t_id = $1 AND u_id = $2", thread.Id, voterId)
	row.Scan(&voted)

	if voice > 0 {

		if voted != 1 {

			voteCounter := 1

			if voted == 0 {
				_, err = Thread.dbLauncher.Exec("INSERT INTO voteThreads (t_id , u_id, counter) VALUES ($1,$2,$3)", thread.Id, voterId, 1)
				voteCounter = 1

			} else {
				_, err = Thread.dbLauncher.Exec("UPDATE voteThreads SET counter = $3 WHERE t_id = $1 AND u_id = $2", thread.Id, voterId, 1)
				voteCounter = 2
			}

			row = Thread.dbLauncher.QueryRow("UPDATE threads SET votes = votes + $2 WHERE t_id = $1 RETURNING date , message, title , votes", thread.Id, voteCounter)
			err = row.Scan(&thread.Created, &thread.Message, &thread.Title, &thread.Votes)

		} else {
			row = Thread.dbLauncher.QueryRow("SELECT date , message, title , votes FROM threads WHERE t_id = $1", thread.Id)
			err = row.Scan(&thread.Created, &thread.Message, &thread.Title, &thread.Votes)
		}
	} else {
		if voted != -1 {

			voteCounter := 0

			if voted == 0 {
				_, err = Thread.dbLauncher.Exec("INSERT INTO voteThreads (t_id , u_id, counter) VALUES ($1,$2, $3)", thread.Id, voterId, -1)
				voteCounter = 1

			} else {
				_, err = Thread.dbLauncher.Exec("UPDATE voteThreads SET counter = $3 WHERE t_id = $1 AND u_id = $2", thread.Id, voterId, -1)
				voteCounter = 2
			}

			row = Thread.dbLauncher.QueryRow("UPDATE threads SET votes = votes - $2 WHERE t_id = $1 RETURNING date , message, title , votes", thread.Id, voteCounter)
			err = row.Scan(&thread.Created, &thread.Message, &thread.Title, &thread.Votes)

		} else {
			row = Thread.dbLauncher.QueryRow("SELECT date , message, title , votes FROM threads WHERE t_id = $1", thread.Id)
			err = row.Scan(&thread.Created, &thread.Message, &thread.Title, &thread.Votes)
		}
	}

	return thread, err

}

func (Thread ThreadRepoRealisation) GetThread(threadId int, thread models.Thread) (models.Thread, error) {

	var row *sql.Row

	if thread.Slug != "" {
		row = Thread.dbLauncher.QueryRow("SELECT TUF.t_id , TUF.slug, U.nickname , F.slug, T.date ,T.message, T.title, T.votes FROM threadUF TUF INNER JOIN users U USING(u_id) INNER JOIN forums F USING(f_id) INNER JOIN threads T USING(t_id) WHERE TUF.slug = $1", thread.Slug)
	} else {
		row = Thread.dbLauncher.QueryRow("SELECT TUF.t_id , TUF.slug, U.nickname , F.slug, T.date ,T.message, T.title, T.votes FROM threadUF TUF INNER JOIN users U USING(u_id) INNER JOIN forums F USING(f_id) INNER JOIN threads T USING(t_id) WHERE TUF.t_id = $1", threadId)
	}

	err := row.Scan(&thread.Id, &thread.Slug, &thread.Author, &thread.Forum, &thread.Created, &thread.Message, &thread.Title, &thread.Votes)

	if err != nil {
		return thread, err
	}

	return thread, nil
}

func (Thread ThreadRepoRealisation) GetPostsSorted(slug string, threadId int, limit int, since int, sortType string, desc bool) ([]models.Message, error) {

	thrdId := 0
	forumSlug := ""
	var row *sql.Row
	if slug != "" {
		row = Thread.dbLauncher.QueryRow("SELECT TUF.t_id , F.slug FROM threadUF TUF INNER JOIN forums F ON(F.f_id=TUF.f_id) WHERE TUF.slug = $1", slug)
	} else {
		row = Thread.dbLauncher.QueryRow("SELECT TUF.t_id , F.slug FROM threadUF TUF INNER JOIN forums F ON(F.f_id=TUF.f_id) WHERE t_id = $1", threadId)
	}

	err := row.Scan(&thrdId, &forumSlug)

	if err != nil {
		return nil, err
	}

	ranger := ">"
	order := "ASC"
	if desc {
		order = "DESC"
		ranger = "<"
	}

	var data *sql.Rows
	messages := make([]models.Message, 0)

	switch sortType {
	case "flat":
		if since != 0 {
			if limit != 0 {
				data, err = Thread.dbLauncher.Query("SELECT M.m_id , M.date, M.message, M.edit, M.parent,U.nickname FROM messages M INNER JOIN messageTU MTU ON(MTU.m_id=M.m_id) INNER JOIN users U ON(U.u_id=MTU.u_id) WHERE MTU.t_id = $1 AND M.m_id"+ranger+"$3 "+" ORDER BY M.m_id "+order+" LIMIT $2", thrdId, limit,since)
			} else {
				data, err = Thread.dbLauncher.Query("SELECT M.m_id , M.date, M.message, M.edit, M.parent,U.nickname FROM messages M INNER JOIN messageTU MTU ON(MTU.m_id=M.m_id) INNER JOIN users U ON(U.u_id=MTU.u_id) WHERE MTU.t_id = $1 AND M.m_id"+ranger+"$2 "+" ORDER BY M.m_id "+order, thrdId,since)
			}
		} else {
			if limit != 0 {
				data, err = Thread.dbLauncher.Query("SELECT M.m_id , M.date, M.message, M.edit, M.parent,U.nickname FROM messages M INNER JOIN messageTU MTU ON(MTU.m_id=M.m_id) INNER JOIN users U ON(U.u_id=MTU.u_id) WHERE MTU.t_id = $1 "+"ORDER BY M.m_id "+order+" LIMIT $2", thrdId, limit)
			} else {
				data, err = Thread.dbLauncher.Query("SELECT M.m_id , M.date, M.message, M.edit, M.parent,U.nickname FROM messages M INNER JOIN messageTU MTU ON(MTU.m_id=M.m_id) INNER JOIN users U ON(U.u_id=MTU.u_id) WHERE MTU.t_id = $1 "+"ORDER BY M.m_id "+order, thrdId)
			}
		}

		if data != nil {

			for data.Next() {
				msg := new(models.Message)
				msg.Forum = forumSlug
				msg.Thread = thrdId
				err = data.Scan(&msg.Id, &msg.Created, &msg.Message, &msg.IsEdited, &msg.Parent, &msg.Author)

				if err != nil {
					fmt.Println(err)
				}

				messages = append(messages, *msg)
			}

			data.Close()
		}

	case "tree":
		if since != 0 {
			if limit != 0 {
				data , err =Thread.dbLauncher.Query("WITH RECURSIVE thread_message(m_id,p_id,date,message,edit,level) AS (SELECT DISTINCT M.m_id,M.parent,M.date,M.message,M.edit,ARRAY[]::BIGINT[] || M.m_id FROM messages M INNER JOIN messageTU MTU ON(M.m_id=MTU.m_id) INNER JOIN users U ON(MTU.u_id=U.u_id AND M.parent=0) WHERE MTU.t_id=$1 UNION ALL SELECT DISTINCT M.m_id, M.parent,M.date,M.message,M.edit,TM.level||M.m_id FROM messages M INNER JOIN thread_message TM ON(TM.m_id=M.parent)) SELECT DISTINCT MT.m_id ,MT.p_id,U.nickname, MT.date,MT.message,MT.edit ,MT.level FROM thread_message MT INNER JOIN messageTU MTU ON(MT.m_id=MTU.m_id) INNER JOIN users U ON (U.u_id=MTU.u_id) WHERE MT.level "+ranger+" (WITH RECURSIVE thread_message(m_id,p_id,date,message,edit,level) AS (SELECT DISTINCT M.m_id,M.parent,M.date,M.message,M.edit,ARRAY[]::BIGINT[] || M.m_id FROM messages M INNER JOIN messageTU MTU ON(M.m_id=MTU.m_id) INNER JOIN users U ON(MTU.u_id=U.u_id AND M.parent=0) WHERE MTU.t_id=$1 UNION ALL SELECT DISTINCT M.m_id, M.parent,M.date,M.message,M.edit,TM.level||M.m_id FROM messages M INNER JOIN thread_message TM ON(TM.m_id=M.parent)) SELECT MT.level FROM thread_message MT WHERE MT.m_id=$3) ORDER BY MT.level "+order+" LIMIT $2" ,thrdId, limit,since)
			} else {
				data , err =Thread.dbLauncher.Query("WITH RECURSIVE thread_message(m_id,p_id,date,message,edit,level) AS (SELECT DISTINCT M.m_id,M.parent,M.date,M.message,M.edit,ARRAY[]::BIGINT[] || M.m_id FROM messages M INNER JOIN messageTU MTU ON(M.m_id=MTU.m_id) INNER JOIN users U ON(MTU.u_id=U.u_id AND M.parent=0) WHERE MTU.t_id=$1 UNION ALL SELECT DISTINCT M.m_id, M.parent,M.date,M.message,M.edit,TM.level||M.m_id FROM messages M INNER JOIN thread_message TM ON(TM.m_id=M.parent)) SELECT DISTINCT MT.m_id ,MT.p_id,U.nickname, MT.date,MT.message,MT.edit ,MT.level FROM thread_message MT INNER JOIN messageTU MTU ON(MT.m_id=MTU.m_id) INNER JOIN users U ON (U.u_id=MTU.u_id) WHERE MT.level "+ranger+" (WITH RECURSIVE thread_message(m_id,p_id,date,message,edit,level) AS (SELECT DISTINCT M.m_id,M.parent,M.date,M.message,M.edit,ARRAY[]::BIGINT[] || M.m_id FROM messages M INNER JOIN messageTU MTU ON(M.m_id=MTU.m_id) INNER JOIN users U ON(MTU.u_id=U.u_id AND M.parent=0) WHERE MTU.t_id=$1 UNION ALL SELECT DISTINCT M.m_id, M.parent,M.date,M.message,M.edit,TM.level||M.m_id FROM messages M INNER JOIN thread_message TM ON(TM.m_id=M.parent)) SELECT MT.level FROM thread_message MT WHERE MT.m_id=$2) ORDER BY MT.level "+order ,thrdId, since)
			}
		} else {
			if limit != 0 {
				data , err =Thread.dbLauncher.Query("WITH RECURSIVE thread_message(m_id,p_id,date,message,edit,level) AS (SELECT DISTINCT M.m_id,M.parent,M.date,M.message,M.edit,ARRAY[]::BIGINT[] || M.m_id FROM messages M INNER JOIN messageTU MTU ON(M.m_id=MTU.m_id AND M.parent=0) INNER JOIN users U ON(MTU.u_id=U.u_id) WHERE MTU.t_id=$1 UNION ALL SELECT DISTINCT M.m_id, M.parent,M.date,M.message,M.edit,TM.level||M.m_id FROM messages M INNER JOIN thread_message TM ON(TM.m_id=M.parent)) SELECT DISTINCT MT.m_id ,MT.p_id,U.nickname, MT.date,MT.message,MT.edit ,MT.level FROM thread_message MT INNER JOIN messageTU MTU ON(MT.m_id=MTU.m_id) INNER JOIN users U ON (U.u_id=MTU.u_id) ORDER BY MT.level "+ order+ " LIMIT $2",thrdId,limit)
			} else {
				data , err = Thread.dbLauncher.Query("WITH RECURSIVE thread_message(m_id,p_id,date,message,edit,level) AS (SELECT DISTINCT M.m_id,M.parent,M.date,M.message,M.edit,ARRAY[]::BIGINT[] || M.m_id FROM messages M INNER JOIN messageTU MTU ON(M.m_id=MTU.m_id AND M.parent=0) INNER JOIN users U ON(MTU.u_id=U.u_id) WHERE MTU.t_id=$1 UNION ALL SELECT DISTINCT M.m_id, M.parent,M.date,M.message,M.edit,TM.level||M.m_id FROM messages M INNER JOIN thread_message TM ON(TM.m_id=M.parent)) SELECT DISTINCT MT.m_id ,MT.p_id,U.nickname, MT.date,MT.message,MT.edit ,MT.level FROM thread_message MT INNER JOIN messageTU MTU ON(MT.m_id=MTU.m_id) INNER JOIN users U ON (U.u_id=MTU.u_id) ORDER BY MT.level" + order,thrdId)
			}
		}

		if err != nil {
			fmt.Println("tree ", err.Error())
		}

		if data != nil {

			for data.Next() {
				msg := new(models.Message)
				msg.Forum = forumSlug
				msg.Thread = thrdId
				level := make([]uint8,0)
				err =data.Scan(&msg.Id , &msg.Parent, &msg.Author, &msg.Created, &msg.Message, &msg.IsEdited, &level)
				if err != nil {
					fmt.Println(err.Error(), "tree")
				}
				messages = append(messages, *msg)
			}

			data.Close()
		}

	case "parent_tree":

		pLevel := 0

		if since != 0 {
			if limit != 0 {

				rwER := make([]uint8,0)

				data , err = Thread.dbLauncher.Query("WITH RECURSIVE thread_message(m_id,p_id,date,message,edit,level) AS (SELECT DISTINCT M.m_id,M.parent,M.date,M.message,M.edit,ARRAY[]::BIGINT[] || M.m_id FROM messages M INNER JOIN messageTU MTU ON(M.m_id=MTU.m_id AND M.parent=0) INNER JOIN users U ON(MTU.u_id=U.u_id) WHERE MTU.t_id=$1 UNION ALL SELECT DISTINCT M.m_id, M.parent,M.date,M.message,M.edit,TM.level||M.m_id FROM messages M INNER JOIN thread_message TM ON(TM.m_id=M.parent) ) SELECT DISTINCT MT.level, MT.level[1] FROM thread_message MT WHERE MT.m_id = $2",thrdId,since)
				data.Next()
				err = data.Scan(&rwER, &pLevel)

				if err != nil {
					fmt.Println("can trash err", err.Error())
				}

				data , err =Thread.dbLauncher.Query("WITH RECURSIVE thread_message(m_id,p_id,date,message,edit,level) AS (SELECT DISTINCT M.m_id,M.parent,M.date,M.message,M.edit,ARRAY[]::BIGINT[] || M.m_id FROM messages M INNER JOIN messageTU MTU ON(M.m_id=MTU.m_id AND M.parent=0) INNER JOIN users U ON(MTU.u_id=U.u_id) WHERE MTU.t_id=$1 UNION ALL SELECT DISTINCT M.m_id, M.parent,M.date,M.message,M.edit,TM.level||M.m_id FROM messages M INNER JOIN thread_message TM ON(TM.m_id=M.parent) ) SELECT DISTINCT MT.m_id ,MT.p_id,U.nickname, MT.date,MT.message,MT.edit ,MT.level,MT.level[1]  FROM thread_message MT INNER JOIN messageTU MTU ON(MT.m_id=MTU.m_id) INNER JOIN users U ON (U.u_id=MTU.u_id) WHERE MT.level "+ranger+" $2 ORDER BY MT.level[1] "+order+" ,MT.level ASC" ,thrdId,rwER)
			} else {
				data , err = Thread.dbLauncher.Query("WITH RECURSIVE thread_message(m_id,p_id,date,message,edit,level) AS (SELECT DISTINCT M.m_id,M.parent,M.date,M.message,M.edit,ARRAY[]::BIGINT[] || M.m_id FROM messages M INNER JOIN messageTU MTU ON(M.m_id=MTU.m_id AND M.parent=0) INNER JOIN users U ON(MTU.u_id=U.u_id) WHERE MTU.t_id=$1 UNION ALL SELECT DISTINCT M.m_id, M.parent,M.date,M.message,M.edit,TM.level||M.m_id FROM messages M INNER JOIN thread_message TM ON(TM.m_id=M.parent) ) SELECT DISTINCT MT.m_id ,MT.p_id,U.nickname, MT.date,MT.message,MT.edit ,MT.level,MT.level[1]  FROM thread_message MT INNER JOIN messageTU MTU ON(MT.m_id=MTU.m_id) INNER JOIN users U ON (U.u_id=MTU.u_id)WHERE MT.level "+ranger+" (WITH RECURSIVE thread_message(m_id,p_id,date,message,edit,level) AS (SELECT DISTINCT M.m_id,M.parent,M.date,M.message,M.edit,ARRAY[]::BIGINT[] || M.m_id FROM messages M INNER JOIN messageTU MTU ON(M.m_id=MTU.m_id AND M.parent=0) INNER JOIN users U ON(MTU.u_id=U.u_id) WHERE MTU.t_id=$1 UNION ALL SELECT DISTINCT M.m_id, M.parent,M.date,M.message,M.edit,TM.level||M.m_id FROM messages M INNER JOIN thread_message TM ON(TM.m_id=M.parent) ) SELECT DISTINCT MT.level FROM thread_message MT WHERE MT.m_id = $2 ORDER BY MT.level "+order+")ORDER BY MT.level[1] "+order+" ,MT.level ASC",thrdId,since)
			}
		} else {
			if limit != 0 {
				data , err =Thread.dbLauncher.Query("WITH RECURSIVE thread_message(m_id,p_id,date,message,edit,level) AS (SELECT DISTINCT M.m_id,M.parent,M.date,M.message,M.edit,ARRAY[]::BIGINT[] || M.m_id FROM messages M INNER JOIN messageTU MTU ON(M.m_id=MTU.m_id AND M.parent=0) INNER JOIN users U ON(MTU.u_id=U.u_id) WHERE MTU.t_id=$1 UNION ALL SELECT DISTINCT M.m_id, M.parent,M.date,M.message,M.edit,TM.level||M.m_id FROM messages M INNER JOIN thread_message TM ON(TM.m_id=M.parent) ) SELECT DISTINCT MT.m_id ,MT.p_id,U.nickname, MT.date,MT.message,MT.edit ,MT.level,MT.level[1] FROM thread_message MT INNER JOIN messageTU MTU ON(MT.m_id=MTU.m_id) INNER JOIN users U ON (U.u_id=MTU.u_id)WHERE(MT.level[1]) IN(WITH RECURSIVE thread_message(m_id,p_id,date,message,edit,level) AS (SELECT DISTINCT M.m_id,M.parent,M.date,M.message,M.edit,ARRAY[]::BIGINT[] || M.m_id FROM messages M INNER JOIN messageTU MTU ON(M.m_id=MTU.m_id AND M.parent=0) INNER JOIN users U ON(MTU.u_id=U.u_id) WHERE MTU.t_id=$1 UNION ALL SELECT DISTINCT M.m_id, M.parent,M.date,M.message,M.edit,TM.level||M.m_id FROM messages M INNER JOIN thread_message TM ON(TM.m_id=M.parent) ) SELECT DISTINCT MT.level[1] FROM thread_message MT ORDER BY MT.level[1] "+order+" LIMIT $2)ORDER BY MT.level[1] "+order+", MT.level ASC",thrdId,limit)
			} else {
				data , err = Thread.dbLauncher.Query("WITH RECURSIVE thread_message(m_id,p_id,date,message,edit,level) AS (SELECT DISTINCT M.m_id,M.parent,M.date,M.message,M.edit,ARRAY[]::BIGINT[] || M.m_id FROM messages M INNER JOIN messageTU MTU ON(M.m_id=MTU.m_id AND M.parent=0) INNER JOIN users U ON(MTU.u_id=U.u_id) WHERE MTU.t_id=$1 UNION ALL SELECT DISTINCT M.m_id, M.parent,M.date,M.message,M.edit,TM.level||M.m_id FROM messages M INNER JOIN thread_message TM ON(TM.m_id=M.parent)) SELECT DISTINCT MT.m_id ,MT.p_id,U.nickname, MT.date,MT.message,MT.edit ,MT.level,MT.level[1] FROM thread_message MT INNER JOIN messageTU MTU ON(MT.m_id=MTU.m_id) INNER JOIN users U ON (U.u_id=MTU.u_id) ORDER BY MT.level[1] " + order+" ,MT.level ASC",thrdId)
			}
		}

		if err != nil {
			fmt.Println("parent tree ", err.Error())
		}

		if data != nil {

			for data.Next() {
				msg := new(models.Message)
				msg.Forum = forumSlug
				msg.Thread = thrdId
				level := make([]uint8,0)
				val := 0
				err =data.Scan(&msg.Id , &msg.Parent, &msg.Author, &msg.Created, &msg.Message, &msg.IsEdited, &level,&val)
				if err != nil {
					fmt.Println(err.Error(), "parent tree")
				}

				if !(since != 0 && pLevel == val) {
					messages = append(messages, *msg)
				}

			}

			data.Close()
		}

	}



	//fmt.Println(messages)
	return messages, nil

}

func (Thread ThreadRepoRealisation) UpdateThread(slug string, threadId int, newThread models.Thread) (models.Thread, error) {

	var row *sql.Row
	if slug != "" {
		row = Thread.dbLauncher.QueryRow("SELECT TUF.t_id , F.slug, U.nickname,TUF.slug FROM threadUF TUF INNER JOIN forums F ON(F.f_id=TUF.f_id) INNER JOIN users U ON(TUF.u_id=U.u_id) WHERE TUF.slug = $1", slug)
	} else {
		row = Thread.dbLauncher.QueryRow("SELECT TUF.t_id , F.slug, U.nickname,TUF.slug FROM threadUF TUF INNER JOIN forums F ON(F.f_id=TUF.f_id) INNER JOIN users U ON(TUF.u_id=U.u_id) WHERE t_id = $1", threadId)
	}

	err := row.Scan(&newThread.Id, &newThread.Forum, &newThread.Author, &newThread.Slug)

	if err != nil {
		return newThread, err
	}

	var threadRow *sql.Rows

	if newThread.Title == "" && newThread.Message == "" {
		threadRow, err = Thread.dbLauncher.Query("SELECT date , message , title , votes FROM threads WHERE t_id = $1", newThread.Id)
		if err != nil {
			fmt.Println(err)
		}

		threadRow.Next()
		err = threadRow.Scan(&newThread.Created, &newThread.Message, &newThread.Title, &newThread.Votes)
		if err != nil {
			fmt.Println(err)
		}

		if threadRow != nil {
			threadRow.Close()
		}

		return newThread, nil
	}

	if newThread.Title != "" && newThread.Message != "" {
		threadRow, err = Thread.dbLauncher.Query("UPDATE threads SET message = $1, title = $2 WHERE t_id = $3 RETURNING date , message, title , votes ", newThread.Message, newThread.Title, newThread.Id)
		if err != nil {
			fmt.Println(err)
		}

		threadRow.Next()
		err = threadRow.Scan(&newThread.Created, &newThread.Message, &newThread.Title, &newThread.Votes)
		if err != nil {
			fmt.Println(err)
		}

		if threadRow != nil {
			threadRow.Close()
		}

		return newThread, nil
	}

	if newThread.Title != "" {
		threadRow, err = Thread.dbLauncher.Query("UPDATE threads SET title = $1 WHERE t_id = $2 RETURNING date , message, title , votes ", newThread.Title, newThread.Id)
		if err != nil {
			fmt.Println(err)
		}

		threadRow.Next()
		err = threadRow.Scan(&newThread.Created, &newThread.Message, &newThread.Title, &newThread.Votes)
		if err != nil {
			fmt.Println(err)
		}

		if threadRow != nil {
			threadRow.Close()
		}

		return newThread, nil
	}

	threadRow, err = Thread.dbLauncher.Query("UPDATE threads SET message = $1 WHERE t_id = $2 RETURNING date , message, title , votes ", newThread.Message, newThread.Id)
	if err != nil {
		fmt.Println(err)
	}

	threadRow.Next()
	err = threadRow.Scan(&newThread.Created, &newThread.Message, &newThread.Title, &newThread.Votes)
	if err != nil {
		fmt.Println(err)
	}

	if threadRow != nil {
		threadRow.Close()
	}

	return newThread, nil

}
