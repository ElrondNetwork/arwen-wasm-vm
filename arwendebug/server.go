package arwendebug

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DebugServer is the debugging server
type DebugServer struct {
	facade  *DebugFacade
	address string
}

// NewDebugServer creates a Server object
func NewDebugServer(facade *DebugFacade, address string) *DebugServer {
	return &DebugServer{
		facade:  facade,
		address: address,
	}
}

// StartServer starts the debugging server
func (server *DebugServer) Start() error {
	log.Debug("Start()")

	router := gin.Default()

	router.POST("/account", server.handleCreateAccount)
	router.POST("/deploy", server.handleDeploy)
	router.POST("/upgrade", server.handleUpgrade)
	router.POST("/run", server.handleRun)
	router.POST("/query", server.handleQuery)

	return router.Run(server.address)
}

func (server *DebugServer) handleCreateAccount(ginContext *gin.Context) {
	request := CreateAccountRequest{}

	err := ginContext.ShouldBindJSON(&request)
	if err != nil {
		returnBadRequest(ginContext, "handleCreateAccount.ShouldBindJSON", err)
		return
	}

	response, err := server.facade.CreateAccount(request)
	if err != nil {
		returnBadRequest(ginContext, "handleCreateAccount.CreateAccount", err)
		return
	}

	returnOkResponse(ginContext, response)
}

func (server *DebugServer) handleDeploy(ginContext *gin.Context) {
	request := DeployRequest{}

	err := ginContext.ShouldBindJSON(&request)
	if err != nil {
		returnBadRequest(ginContext, "handleDeploy.ShouldBindJSON", err)
		return
	}

	response, err := server.facade.DeploySmartContract(request)
	if err != nil {
		returnBadRequest(ginContext, "handleDeploy.DeploySmartContract", err)
		return
	}

	returnOkResponse(ginContext, response)
}

func (server *DebugServer) handleUpgrade(ginContext *gin.Context) {
	request := UpgradeRequest{}

	err := ginContext.ShouldBindJSON(&request)
	if err != nil {
		returnBadRequest(ginContext, "handleUpgrade.ShouldBindJSON", err)
		return
	}

	response, err := server.facade.UpgradeSmartContract(request)
	if err != nil {
		returnBadRequest(ginContext, "handleUpgrade.UpgradeSmartContract", err)
		return
	}

	returnOkResponse(ginContext, response)
}

func (server *DebugServer) handleRun(ginContext *gin.Context) {
	request := RunRequest{}

	err := ginContext.ShouldBindJSON(&request)
	if err != nil {
		returnBadRequest(ginContext, "handleRun.ShouldBindJSON", err)
		return
	}

	response, err := server.facade.CallSmartContract(request)
	if err != nil {
		returnBadRequest(ginContext, "handleRun.UpgradeSmartContract", err)
		return
	}

	returnOkResponse(ginContext, response)
}

func (server *DebugServer) handleQuery(ginContext *gin.Context) {
	request := QueryRequest{}

	err := ginContext.ShouldBindJSON(&request)
	if err != nil {
		returnBadRequest(ginContext, "handleQuery.ShouldBindJSON", err)
		return
	}

	response, err := server.facade.QuerySmartContract(request)
	if err != nil {
		returnBadRequest(ginContext, "handleQuery.UpgradeSmartContract", err)
		return
	}

	returnOkResponse(ginContext, response)
}

func returnBadRequest(context *gin.Context, errScope string, err error) {
	context.JSON(http.StatusBadRequest, gin.H{
		"error":        fmt.Sprintf("%T", err),
		"errorMessage": err.Error(),
		"errorScope":   errScope,
		"data":         nil,
	})
}

func returnOkResponse(context *gin.Context, data interface{}) {
	context.JSON(http.StatusOK, gin.H{
		"error":        nil,
		"errorMessage": nil,
		"errorScope":   nil,
		"data":         data,
	})
}
