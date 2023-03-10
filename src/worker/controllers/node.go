package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sonntuet1997/avalanche-simplified/worker/services"
	baseEx "gitlab.com/golibs-starter/golib/exception"
	"gitlab.com/golibs-starter/golib/log"
	"gitlab.com/golibs-starter/golib/web/response"
	"strconv"
	"strings"
)

type NodeController struct {
	ConsensusService *services.ConsensusService
}

func NewNodeController(
	ConsensusService *services.ConsensusService,
) *NodeController {
	return &NodeController{
		ConsensusService: ConsensusService,
	}
}

// GetConfirmedTransactions ...
//
//	@Summary		Get Confirmed Transactions
//	@Description	Return confirmed transactions
//	@Tags			Node
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response{data=[]entities.Transaction}
//	@Failure		500	{object}	response.Response
//	@Router			/v1/node/transactions [get]
func (r *NodeController) GetConfirmedTransactions(ctx *gin.Context) {
	result := r.ConsensusService.GetConfirmedTransactions()
	response.Write(ctx.Writer, response.Ok(result))
}

// GetPreferTransaction ...
//
//	@Summary		Get Prefer Transaction
//	@Description	Return prefer transaction
//	@Tags			Node
//	@Accept			json
//	@Produce		json
//	@Param			block_number	path		string	true	"block number"
//	@Success		200				{object}	response.Response{data=entities.Transaction}
//	@Failure		500				{object}	response.Response
//	@Router			/v1/node/prefer-transactions/{block_number} [get]
func (r *NodeController) GetPreferTransaction(ctx *gin.Context) {
	blockNumberString := strings.TrimSpace(ctx.Param("block_number"))
	blockNumber, err := strconv.Atoi(blockNumberString)
	if err != nil {
		response.WriteError(ctx.Writer, baseEx.BadRequest)
		return
	}
	preference, err := r.ConsensusService.MyPreference(blockNumber)
	if err != nil {
		log.Error(ctx, "[NodeController][GetPreferTransaction] get MyPreference %+v failed with err: %+v", blockNumber, err)
		response.WriteError(ctx.Writer, baseEx.SystemError)
		return
	}
	response.Write(ctx.Writer, response.Ok(preference))
}
