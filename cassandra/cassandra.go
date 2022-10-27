package cassandra

import (
	"fmt"
	"log"
	"time"

	models "github.com/HarishGangula/content-state-go-concurrency/models"
	"github.com/gocql/gocql"
)

var session *gocql.Session

func Init() {
	var err error
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "sunbirdcourse"
	session, err = cluster.CreateSession()

	if err != nil {
		panic(err)
	}
	fmt.Println("Cassandra init is done")

	tableCreateStatement := "CREATE TABLE IF NOT EXISTS contentstate (userId text, contentId text, batchId text, courseId text, status int, lastAccessTime text,    PRIMARY KEY (contentId, batchId, courseId, userId))"
	if tableErr := session.Query(tableCreateStatement).Exec(); tableErr != nil {
		panic(tableErr)
	} else {
		fmt.Println("content state table created successfully")
	}

}

func Close() {
	session.Close()
}

func UpsertContentState(r models.Request) map[string]bool {
	response := map[string]bool{}
	for _, content := range r.Request.Contents {
		// session.Query()
		query := `UPDATE contentstate  SET status=?, lastAccessTime=? WHERE userId=? and contentId=? and courseId=? and batchId=?`
		if err := session.Query(query, content.Status, string(time.Now().UnixMilli()), r.Request.UserId, content.ContentId, content.CourseId, content.BatchId).Exec(); err != nil {
			log.Println("Update error:", err)
			response[content.ContentId] = false
		} else {
			response[content.ContentId] = true
		}
	}
	return response
}
