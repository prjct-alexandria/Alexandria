package postgres

import (
	"database/sql"
	"fmt"
	"mainServer/entities"
	"mainServer/models"
)

type PgCommitThreadRepository struct {
	Db *sql.DB
}

// GetCommitThreads  gets the commit comment threads from the database, using the article id (aid) and commit id (cid)
func (r PgCommitThreadRepository) GetCommitThreads(aid int64, cid int64) ([]models.Thread, error) {
	// construct list of threads
	stmt, err := r.Db.Prepare("SELECT t.threadid, articleid, commitid, commentid, authorid, creationdate, content FROM committhread ct JOIN thread t on ct.threadid = t.threadid JOIN comment c on t.threadid = c.threadid WHERE t.articleid = $1 AND ct.commitid = $2")
	if err != nil {
		return nil, fmt.Errorf("GetCommitThreads: %v", err)
	}

	rows, err := stmt.Query(aid, cid)
	if err != nil {
		return nil, fmt.Errorf("GetCommitThreads: %v", err)
	}
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("GetCommitThreads: %v", err)
	}

	var threads []models.Thread
	for rows.Next() {
		// declare variables
		var tid int64
		var aid int64
		var cid int64
		var coid int64
		var auid string
		var cd string
		var c string
		if err := rows.Scan(&tid, &aid, &cid, &coid, &auid, &cd, &c); err != nil {
			return threads, fmt.Errorf("GetCommitThreads: %v", err)
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
				SpecificId: cid,
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
		return threads, fmt.Errorf("GetCommitThreads: %v", err)
	}
	return threads, nil
}

func NewPgCommitThreadRepository(db *sql.DB) PgCommitThreadRepository {
	repo := PgCommitThreadRepository{db}
	err := repo.createCommitThreadTable()
	if err != nil {
		return PgCommitThreadRepository{}
	}
	return repo
}

func (r PgCommitThreadRepository) createCommitThreadTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS commitThread (
    	commitThreadId SERIAL,
    	commitId BIGINT NOT NULL,
        threadId BIGINT NOT NULL,
    	PRIMARY KEY (commitThreadId)
    )`)
	return err
}

func (r PgCommitThreadRepository) CreateCommitThread(thread models.Thread, tid int64) (int64, error) {
	stmt, err := r.Db.Prepare(`INSERT INTO committhread (committhreadid, commitid, threadid) 
	VALUES (DEFAULT, $1, $2) RETURNING committhreadid`)
	if err != nil {
		return -1, fmt.Errorf("CreateThread: %v", err)
	}
	row := stmt.QueryRow(thread.SpecificId, tid)

	var id int64
	err = row.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("CreateThread: %v", err)
	}

	return id, nil
}
