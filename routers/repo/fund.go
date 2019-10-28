// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repo

import (
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
	ctx.Data["IsRepoToolbarFunds"] = true
	page := ctx.QueryInt("page")
	funds, err := models.GetFunds(ctx.Repo.Repository.ID, page)
	if err != nil {
		ctx.ServerError("GetFunds", err)
		return
	}
	ctx.Data["Funds"] = funds
	ctx.HTML(200, tplFunds)
}

