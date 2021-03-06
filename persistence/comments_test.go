package persistence

import (
	"database/sql"
	"reflect"
	"testing"
)

func TestComments(t *testing.T) {
	tests := commentTests{db}

	data := tests.prepareCommentTestsInput()

	t.Run("get", tests.get(data.comments[0]))
	t.Run("list", tests.list(data.task.ID))
	t.Run("add", tests.add(data.task.ID))
	t.Run("upd", tests.upd(data.comments[0]))
	t.Run("del", tests.del(data.comments[0]))

}

type commentTests struct {
	*sql.DB
}
type commentTestsInput struct {
	task     Task
	comments []Comment
}

func (test *commentTests) prepareCommentTestsInput() commentTestsInput {
	projects, err := prepareProjects(test.DB, 1)
	panicOnErr(err)
	statuses, err := prepareStatuses(test.DB, projects[0], 1)
	panicOnErr(err)
	priorities, err := preparePriorities(test.DB, 1)
	panicOnErr(err)
	tasks, err := prepareTasks(test.DB, projects[0], statuses[0], priorities[0], 1)
	panicOnErr(err)
	comments := prepareComments(test.DB, tasks[0].ID, 2)
	return commentTestsInput{
		task:     tasks[0],
		comments: comments,
	}
}

func (test *commentTests) add(taskID int64) func(t *testing.T) {
	return func(t *testing.T) {
		c := Comment{
			TaskID: taskID,
			Text:   "testc",
		}
		_, err := AddComment(c)
		if err != nil {
			t.Error(err)
		}
	}
}

func (test *commentTests) get(expected Comment) func(t *testing.T) {
	return func(t *testing.T) {
		actual, err := GetComment(expected.ID)
		if err != nil {
			t.Fatal(err)
		}
		actual.Date = expected.Date
		if !reflect.DeepEqual(expected, actual) {
			t.Error("expected/actual mismatch:", expected, actual)
		}
	}
}

func (test *commentTests) list(taskID int64) func(t *testing.T) {
	return func(t *testing.T) {
		comments, err := ListComments(taskID)
		if err != nil {
			t.Fatal(err)
		}
		for _, c := range comments {
			if !checkCommentExists(test.DB, c) {
				t.Error("missing:", c)
			}
		}
	}
}

func (test *commentTests) upd(old Comment) func(t *testing.T) {
	upd := old
	upd.Text += "_updated"
	return func(t *testing.T) {
		err := UpdComment(upd)
		if err != nil {
			t.Fatal(err)
		}
		if !checkCommentExists(test.DB, upd) {
			t.Error("updated comment missing")
		}
	}
}

func (test *commentTests) del(c Comment) func(t *testing.T) {
	return func(t *testing.T) {
		err := DelComment(c.ID)
		if err != nil {
			t.Fatal(err)
		}
		if checkCommentExists(test.DB, c) {
			t.Error("not deleted")
		}
	}
}

func checkCommentExists(db *sql.DB, expected Comment) bool {
	q := `SELECT * FROM comments WHERE id=$1`
	var actual Comment
	err := db.QueryRow(q, expected.ID).Scan(&actual.ID, &actual.TaskID, &actual.Text, &actual.Date)
	actual.Date = expected.Date
	return err == nil && reflect.DeepEqual(expected, actual)
}
