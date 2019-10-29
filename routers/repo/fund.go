// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repo

import (
	"strconv"
	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/context"
)

const (
	tplFunds base.TplName = "repo/fund/list"
)

// Funds render repository branch page
func Funds(ctx *context.Context) {
	ctx.Data["Title"] = "Funds"
	ctx.Data["PageIsFund"] = true
	page := ctx.QueryInt("page")
	funds, err := models.GetFunds(ctx.Repo.Repository.ID, page)
	if err != nil {
		ctx.ServerError("GetFunds", err)
		return
	}
	ctx.Data["Funds"] = funds
	ctx.HTML(200, tplFunds)
}

func Funding(ctx *context.Context) {
	//减少发送者点数
	//增加接收者点数
	//创建transfer记录
	var err error
	Qty, err := strconv.Atoi(ctx.Query("qty"))
	if err != nil {
		ctx.Flash.Error("请输入数字")
		return
	}
	var repoid int
	repoid, err = strconv.Atoi(ctx.Query("repoid"))
	if ctx.User.Point < Qty {
		ctx.Flash.Error("余额不足！")
		return
	}
	err = models.TransferPoint(ctx.User.Name,
			ctx.Query("why"),
			ctx.Query("toid"),
			Qty)
	if err != nil {
	ctx.ServerError("Transfer", err)
		return
	}
	err = models.NewFund(ctx.User.Name,
		int64(repoid),
		int64(Qty))
	if err != nil {
		ctx.ServerError("Transfer", err)
		return
	}

	ctx.RedirectToFirst(ctx.Query("redirect_to"), ctx.Repo.RepoLink + "/funds")
}

