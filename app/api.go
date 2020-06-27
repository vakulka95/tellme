package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"gitlab.com/tellmecomua/tellme.api/app/config"
	"gitlab.com/tellmecomua/tellme.api/app/persistence"
	"gitlab.com/tellmecomua/tellme.api/app/persistence/pgx"
	"gitlab.com/tellmecomua/tellme.api/pkg/sms"
	"gitlab.com/tellmecomua/tellme.api/pkg/util/job"
)

type apiserver struct {
	engine     *gin.Engine
	config     *config.Environment
	smscli     *sms.TurboSMS
	repository persistence.Repository

	documentHostDir   string
	documentServePath string
	swaggerServeFile  string
}

func New(c *config.Environment) Application {
	return &apiserver{config: c}
}

func (s *apiserver) init() error {

	s.documentHostDir = s.config.ExpertDocumentsStoreDir
	s.documentServePath = "document/expert"

	s.registerHandlers()
	s.smscli = sms.NewManager(s.config.TurboSMSUsername, s.config.TurboSMSPassword)
	return s.connectRepository()
}

func (s *apiserver) connectRepository() error {
	s.repository = pgx.New(s.config)
	return s.repository.Connect()
}

func (s *apiserver) Run() error {
	if err := s.init(); err != nil {
		return fmt.Errorf("failed to init apiserver: %v", err)
	}
	defer func() {
		if err := s.repository.Disconnect(); err != nil {
			log.Printf("failed to close db connection: %v", err)
		}
	}()

	addr := fmt.Sprintf(":%s", s.config.ServePort)
	handler := &http.Server{
		Handler: s.engine,
		Addr:    addr,
	}

	go func() {
		log.Printf("Start listening and serving [http] on %s...", addr)
		if err := handler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http listen and serve: %v", err)
		}
	}()

	notProcessedRequisitionReplyJob := job.Run(s.config.NotProcessedRequisitionJobPeriod, "NotProcessedRequisitionReply", s.jobNotProcessedRequisitionReply)
	defer notProcessedRequisitionReplyJob.Stop()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shuting down server...")
	return handler.Shutdown(ctx)
}

func (s *apiserver) generateDocumentURL(expertID, filename string) string {
	return fmt.Sprintf("https://%s/%s/%s/%s", s.config.DomainName, s.documentServePath, expertID, filename)
}
