package week2

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// 第二周作业
// 我认为dao层中遇到sql.ErrNoRows时，不用Wrap该error，抛给上层。
// 因为sql.ErrNoRows表示"sql: no rows in result set"，即没有查找到相应的行记录。
// 我的做法是将需要Scan的对象置为nil，并将err也置为nil，然后return，这样上层得到值为nil的对象时，
// 也可以知道没有查询到相应的行记录。
// 而对于其他错误情况，则对err进行Wrap，然后return

var (
	db *sql.DB
)

// Init 初始化mysql连接
func Init() {
	sqlDB, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/practice")
	if err != nil {
		log.Fatal(err)
	}
	db = sqlDB
}

// Note 表示笔记, 对应于数据库中的note表
type Note struct {
	ID int64 `json:"id"`
	UserID int64 `json:"user_id"`
	Title string `json:"title"`
	Content string `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (n *Note) TableName() string {
	return "note"
}

// QueryNoteById 按照note的id查询note
func QueryNoteById(id int64) (*Note, error) {
	var (
		// 查询一行
		row = db.QueryRow("SELECT id, user_id, title, content FROM note WHERE id=?", id)
		n = &Note{}
		err error
	)
	if err = row.Scan(&n.ID, &n.UserID, &n.Title, &n.Content); err != nil {
		// 遇到sql.ErrNoRows表示没有查询到行记录
		if err == sql.ErrNoRows {
			n = nil
			err = nil
		} else {
			// 其他错误情况则对err进行Wrap
			n = nil
			err = fmt.Errorf("row.Scan err: %w", err)
		}
		return n, err
	}
	// 正常返回
	return n, err
}