// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"fmt"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/modules/timeutil"
	"xorm.io/builder"
)

type Transfer struct {
	ID       int64 `xorm:"pk autoincr"`
	FromID   int64 
	ToID     int64
	Why      string
	Qty      int
	CreatedUnix   timeutil.TimeStamp `xorm:"INDEX created"`
}

// SearchTransOptions contains the options for searching
type SearchTransOptions struct {
	Keyword       string
	OrderBy       SearchOrderBy
	Page          int
	PageSize      int   // Can be smaller than or equal to setting.UI.ExplorePagingNum
}

func (opts *SearchTransOptions) toConds() builder.Cond {
	var cond builder.Cond = builder.Gt{"ID": 0}
	return cond
}

func SearchTrans(opts *SearchTransOptions) (trans []*Transfer, _ int64, _ error) {
	cond := opts.toConds()
	count, err := x.Where(cond).Count(new(Transfer))
	if err != nil {
		return nil, 0, fmt.Errorf("Count: %v", err)
	}

	if opts.PageSize == 0 || opts.PageSize > setting.UI.ExplorePagingNum {
		opts.PageSize = setting.UI.ExplorePagingNum
	}
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if len(opts.OrderBy) == 0 {
		opts.OrderBy = SearchOrderByIDReverse
	}

	sess := x.Where(cond)
	if opts.PageSize > 0 {
		sess = sess.Limit(opts.PageSize, (opts.Page-1)*opts.PageSize)
	}
	if opts.PageSize == -1 {
		opts.PageSize = int(count)
	}

	trans = make([]*Transfer, 0, opts.PageSize)
	return trans, count, sess.OrderBy(opts.OrderBy.String()).Find(&trans)
}

func TransferPoint(FromID int64, Why string, ToID int64, Qty int) (err error) {

	sess := x.NewSession()

	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return err
	}

	if _, err = sess.Insert(&Transfer{FromID: FromID, ToID: ToID, Why: Why, Qty: Qty}); err != nil {
		return err
	}

	if _, err = sess.Exec("UPDATE `user` SET point = point + ? WHERE id = ?", Qty, ToID); err != nil {
		return err
	}

	if _, err = sess.Exec("UPDATE `user` SET point = point - ? WHERE id = ?", Qty, FromID); err != nil {
		return err
	}
	return sess.Commit()
}

// Fund contains the fund information
type Fund struct {
	ID                int64 `xorm:"pk autoincr"`
	Name              string
	RepoID            int64 `xorm:"INDEX"`
	Qty               int64
}

type FundList []*Fund

// GetFunds returns a list of Funds of given repository.
func GetFunds(repoID int64, page int) (FundList, error) {
	funds := make([]*Fund, 0, setting.UI.IssuePagingNum)
	sess := x.Where("repo_id = ? ", repoID)
	if page > 0 {
		sess = sess.Limit(setting.UI.IssuePagingNum, (page-1)*setting.UI.IssuePagingNum)
	}

	return funds, sess.Find(&funds)
}
