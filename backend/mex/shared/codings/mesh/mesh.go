package mesh

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	// Sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/psanford/sqlite3vfs"

	"github.com/d4l-data4life/mex/mex/shared/codings"
)

type meshSet struct {
	db *sql.DB

	pstmts []*sql.Stmt
}

//nolint:gomnd
func NewCodingsetBytes(sqliteFileContents []byte) (codings.Codingset, error) {
	err := sqlite3vfs.RegisterVFS(VFSName, NewReadOnlyInMemoryVFS(sqliteFileContents))
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("ignore.db?vfs=%s", VFSName))
	if err != nil {
		return nil, err
	}

	ms := &meshSet{
		db:     db,
		pstmts: make([]*sql.Stmt, 10),
	}

	for i := range ms.pstmts {
		ms.pstmts[i], err = db.Prepare(fmt.Sprintf("SELECT Term FROM Terms WHERE DUID IN (%s) AND language = ? AND Type = 6", qmarks(i+1)))
		if err != nil {
			panic(err)
		}
	}

	return ms, nil
}

func NewCodingsetFile(sqliteFile *os.File) (codings.Codingset, error) {
	db, err := sql.Open("sqlite3", sqliteFile.Name())
	if err != nil {
		return nil, err
	}

	return &meshSet{db: db}, nil
}

func (ms *meshSet) Close() {
	for i := range ms.pstmts {
		_ = ms.pstmts[i].Close()
	}
	_ = ms.db.Close()
}

func (ms *meshSet) Info() (map[string]string, error) {
	rows, err := ms.db.Query("SELECT Key, Value FROM Info")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	info := map[string]string{}
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		info[k] = v
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return info, nil
}

func (ms *meshSet) Count() int {
	row := ms.db.QueryRow("SELECT count(DISTINCT DUID) FROM Terms")
	if row == nil {
		return -1
	}

	var v int
	err := row.Scan(&v)
	if err != nil {
		return -1
	}
	return v
}

func (ms *meshSet) GetMainHeadings() ([]string, error) {
	rows, err := ms.db.Query("SELECT DISTINCT DUID FROM Terms WHERE Type = 1 OR Type = 0 ORDER BY DUID ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var duids []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		duids = append(duids, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return duids, nil
}

//nolint:gomnd,lll
func (ms *meshSet) ResolveMainHeadings(descriptorIDs []string, language string) ([]string, error) {
	var err error
	var rows *sql.Rows

	// Test showed improvements when using prepared statements for different sizes of code numbers.
	// We need to verify those in production and then decide whether to keep the below code or revert it back.
	switch len(descriptorIDs) {
	case 0:
		return []string{}, nil
	case 1:
		rows, err = ms.pstmts[0].Query(descriptorIDs[0], language)
	case 2:
		rows, err = ms.pstmts[1].Query(descriptorIDs[0], descriptorIDs[1], language)
	case 3:
		rows, err = ms.pstmts[2].Query(descriptorIDs[0], descriptorIDs[1], descriptorIDs[2], language)
	case 4:
		rows, err = ms.pstmts[3].Query(descriptorIDs[0], descriptorIDs[1], descriptorIDs[2], descriptorIDs[3], language)
	case 5:
		rows, err = ms.pstmts[4].Query(descriptorIDs[0], descriptorIDs[1], descriptorIDs[2], descriptorIDs[3], descriptorIDs[4], language)
	case 6:
		rows, err = ms.pstmts[5].Query(descriptorIDs[0], descriptorIDs[1], descriptorIDs[2], descriptorIDs[3], descriptorIDs[4], descriptorIDs[5], language)
	case 7:
		rows, err = ms.pstmts[6].Query(descriptorIDs[0], descriptorIDs[1], descriptorIDs[2], descriptorIDs[3], descriptorIDs[4], descriptorIDs[5], descriptorIDs[6], language)
	case 8:
		rows, err = ms.pstmts[7].Query(descriptorIDs[0], descriptorIDs[1], descriptorIDs[2], descriptorIDs[3], descriptorIDs[4], descriptorIDs[5], descriptorIDs[6], descriptorIDs[7], language)
	case 9:
		rows, err = ms.pstmts[8].Query(descriptorIDs[0], descriptorIDs[1], descriptorIDs[2], descriptorIDs[3], descriptorIDs[4], descriptorIDs[5], descriptorIDs[6], descriptorIDs[7], descriptorIDs[8], language)
	case 10:
		rows, err = ms.pstmts[9].Query(descriptorIDs[0], descriptorIDs[1], descriptorIDs[2], descriptorIDs[3], descriptorIDs[4], descriptorIDs[5], descriptorIDs[6], descriptorIDs[7], descriptorIDs[8], descriptorIDs[9], language)
	default:
		rows, err = ms.db.Query(fmt.Sprintf("SELECT Term FROM Terms WHERE DUID IN (%s) AND language = $1", quote(descriptorIDs)), language)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var terms []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		terms = append(terms, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return terms, nil
}

func (ms *meshSet) ResolveTreeNumbers(descriptorIDs []string) ([]string, error) {
	rows, err := ms.db.Query(fmt.Sprintf("SELECT Term FROM Terms WHERE DUID IN (%s) AND Type = 2", quote(descriptorIDs)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tnums []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		tnums = append(tnums, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tnums, nil
}

func quote(ss []string) string {
	s := make([]string, len(ss))
	for i := range ss {
		s[i] = fmt.Sprintf("'%s'", ss[i])
	}
	return strings.Join(s, ",")
}

func qmarks(n int) string {
	q := make([]string, n)
	for i := range q {
		q[i] = "?"
	}
	return strings.Join(q, ", ")
}
