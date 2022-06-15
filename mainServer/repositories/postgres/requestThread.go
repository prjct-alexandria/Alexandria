package postgres

import (
	"database/sql"
	"fmt"
	"mainServer/entities"
	"mainServer/models"
)

type PgRequestThreadRepository struct {
	Db *sql.DB
}

//GetRequestThreads gets the request comment threads from the database, using the article id (aid) and request id (rid)
func (r PgRequestThreadRepository) GetRequestThreads(aid int64, rid int64) ([]models.Thread, error) {
	// construct list of threads
	stmt, err := r.Db.Prepare("SELECT t.threadid, articleid, requestid, commentid, authorid, creationdate, content FROM requestthread rt JOIN thread t on rt.threadid = t.threadid JOIN comment c on t.threadid = c.threadid WHERE t.articleid = $1 AND rt.requestid = $2 ORDER BY creationdate")
	if err != nil {
		return nil, fmt.Errorf("GetRequestThreads: %v", err)
	}

	rows, err := stmt.Query(aid, rid)
	if err != nil {
		return nil, fmt.Errorf("GetRequestThreads: %v", err)
	}
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("GetRequestThreads: %v", err)
	}

	var threads []models.Thread
	for rows.Next() {
		// declare variables
		var tid int64
		var aid int64
		var rid int64
		var coid int64
		var auid string
		var cd string
		var c string
		if err := rows.Scan(&tid, &aid, &rid, &coid, &auid, &cd, &c); err != nil {
			return threads, fmt.Errorf("GetRequestThreads: %v", err)
		}
		// find index of tid in the threads list
		index := -1
		for i, v := range threads {
			if v.Id == tid {
				index = i
				break
			}
		}
		// construct comment
		comment := entities.Comment{
			Id:           coid,
			AuthorId:     auid,
			ThreadId:     tid,
			Content:      c,
			CreationDate: cd,
		}
		// if index tid isn't yet in the list, add new thread, else add comment to thread
		if index == -1 {
			var comments []entities.Comment
			comments = append(comments, comment)
			threads = append(threads, models.Thread{
				Id:         tid,
				ArticleId:  aid,
				SpecificId: rid,
				Comment:    comments,
			})
		} else {
			oldThread := threads[index]
			oldComments := oldThread.Comment
			newComments := append(oldComments, comment)
			threads[index] = models.Thread{
				Id:         oldThread.Id,
				ArticleId:  oldThread.ArticleId,
				SpecificId: oldThread.SpecificId,
				Comment:    newComments,
			}
		}
	}
	if err = rows.Err(); err != nil {
		return threads, fmt.Errorf("GetRequestThreads: %v", err)
	}
	return threads, nil
}

func NewPgRequestThreadRepository(db *sql.DB) PgRequestThreadRepository {
	repo := PgRequestThreadRepository{db}
	err := repo.createRequestThreadTable()
	if err != nil {
		return PgRequestThreadRepository{}
	}
	return repo
}

func (r PgRequestThreadRepository) createRequestThreadTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS RequestThread (
    	RequestThreadId SERIAL,
    	RequestId BIGINT NOT NULL,
        threadId BIGINT NOT NULL,
    	PRIMARY KEY (RequestThreadId)
    )`)
	return err
}

func (r PgRequestThreadRepository) CreateRequestThread(thread models.Thread, tid int64) (int64, error) {
	stmt, err := r.Db.Prepare(`INSERT INTO requestthread (requestthreadid, requestid, threadid)
	VALUES (DEFAULT, $1, $2) RETURNING requestthreadid`)
	if err != nil {
		return -1, fmt.Errorf("CreateThread: %v", err)
	}
	row := stmt.QueryRow(thread.SpecificId, tid)

	var id int64
	err = row.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("CreateThread: %v", err)
	}

	return id, err
}
