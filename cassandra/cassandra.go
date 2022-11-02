package cassandra

import (
	"fmt"
	"log"
	"sync"
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

func update(userId string, content models.Content, res chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	query := `UPDATE contentstate  SET status=?, lastAccessTime=? WHERE userId=? and contentId=? and courseId=? and batchId=?`
	if err := session.Query(query, content.Status, string(time.Now().UnixMilli()), userId, content.ContentId, content.CourseId, content.BatchId).Exec(); err != nil {
		log.Println("Update error:", err)
		res <- false
	} else {
		res <- true
	}
}

func UpsertContentState(r models.Request, responseChan chan models.Response) {
	wg := new(sync.WaitGroup)
	response := models.Response{}
	channel := make(chan bool)
	for _, content := range r.Request.Contents {
		wg.Add(1)
		go update(r.Request.UserId, content, channel, wg)
		isUpserted := <-channel
		if isUpserted {
			response[content.ContentId] = true
		} else {
			response[content.ContentId] = false
		}

	}
	go func() {
		wg.Wait()
		defer close(channel)
	}()

	responseChan <- response
	close(responseChan)
}
