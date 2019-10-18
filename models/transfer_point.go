// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"code.gitea.io/gitea/modules/timeutil"
)

type Transfer struct {
	ID       int64 `xorm:"pk autoincr"`
	FromID   int64 
	ToID     int64
	Why      string
	Qty      int
	CreatedUnix   timeutil.TimeStamp `xorm:"INDEX created"`
}

func TransferPoint(FromID int64, Why string, ToID int64, Qty int) (err error) {

	sess := x.NewSession()
	//判断是否足够转
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
