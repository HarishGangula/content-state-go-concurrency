const express = require('express')
const app = express()
app.use(express.json())
app.patch('/api/v1/contentstate/update', async (req, res) => {
  await update(req.body)
  res.send({"success": "ok"})
})

const cassandra = require('cassandra-driver');

const client = new cassandra.Client({
  contactPoints: ['127.0.0.1'],
  localDataCenter: "datacenter1",
  keyspace: 'sunbirdcourse'
});
const update = async (body) => {
    const userId = body?.request?.userId;
    const contents = body?.request?.contents;
    for (const content of contents) {
        await client.execute(
            "UPDATE contentstate  SET status=?, lastAccessTime=? WHERE userId=? and contentId=? and courseId=? and batchId=?"
            , [content.status, Date.now()+"", userId, content.contentId, content.courseId, content.batchId], { prepare: true });
    }
}

app.listen(3000)