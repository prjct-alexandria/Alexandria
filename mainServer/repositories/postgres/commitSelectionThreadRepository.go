package postgres

import (
	"database/sql"
	"fmt"
	"mainServer/entities"
	"mainServer/models"
)

type PgCommitSectionThreadRepository struct {
	Db *sql.DB
}

func NewPgCommitSelectionThreadRepository(db *sql.DB) PgCommitSectionThreadRepository {
	repo := PgCommitSectionThreadRepository{db}
	err := repo.createCommitSelectionThreadTable()
	if err != nil {
		return PgCommitSectionThreadRepository{}
	}
	return repo
}

func (r PgCommitSectionThreadRepository) createCommitSelectionThreadTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS commitselectionthread (
    	commitSectionThreadId SERIAL,
    	commitId NCHAR(40) NOT NULL,
        threadId BIGINT NOT NULL,
        section  NCHAR(255) NOT NULL,
    	PRIMARY KEY (commitSectionThreadId)
    )`)
	return err
}

func (r PgCommitSectionThreadRepository) CreateCommitSelectionThread(cid string, tid int64, section string) (int64, error) {
	stmt, err := r.Db.Prepare(`INSERT INTO commitselectionthread (commitsectionthreadid, commitid, threadid, section)
	VALUES (DEFAULT, $1, $2, $3) RETURNING commitsectionthreadid`)
	if err != nil {
		return -1, fmt.Errorf("CreateThread: %v", err)
	}
	row := stmt.QueryRow(cid, tid, section)

	var id int64
	err = row.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("CreateThread: %v", err)
	}

	return id, nil
}

func (r PgCommitSectionThreadRepository) GetCommitSelectionThreads(aid int64, cid string) ([]models.SelectionThread, error) {
	// construct list of threads
	stmt, err := r.Db.Prepare(`SELECT t.threadid, articleid, commitid, commentid, authorid, creationdate, content, section
	FROM commitselectionthread ct JOIN thread t on ct.threadid = t.threadid JOIN comment c on t.threadid = c.threadid
	WHERE t.articleid = $1 AND ct.commitid = $2 ORDER BY creationdate`)

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

	var threads []models.SelectionThread
	for rows.Next() {
		// declare variables
		var tid int64
		var aid int64
		var cid string
		var coid int64
		var auid string
		var cd string
		var c string
		var section string
		if err := rows.Scan(&tid, &aid, &cid, &coid, &auid, &cd, &c, &section); err != nil {
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
			threads = append(threads, models.SelectionThread{
				Id:         tid,
				ArticleId:  aid,
				SpecificId: cid,
				Comments:   comments,
				Selection:  section,
			})
		} else {
			threads[index].Comments = append(threads[index].Comments, comment)
		}
	}
	if err = rows.Err(); err != nil {
		return threads, fmt.Errorf("GetCommitThreads: %v", err)
	}
	return threads, nil
}
