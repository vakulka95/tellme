package app

import (
	"fmt"
	"log"
	"time"
)

func (s *apiserver) jobNotProcessedRequisitionReply() error {
	lifetimeTo := time.Now().Add(-48 * time.Hour) // 48 hours

	requisitions, err := s.repository.GetNotProcessedRequisition(&lifetimeTo)
	if err != nil {
		return fmt.Errorf("failed to get not processed requisitions: %v", err)
	}

	for _, requisition := range requisitions {

		if requisition.SMSReplyCount != 0 { // means already processed
			continue
		}

		err = s.smscli.SendRequisitionReply(requisition.Phone, requisition.Username, requisition.Gender)

		if err != nil {
			log.Printf("(ERR) Failed to send requisiton reply: %v", err)
			continue
		}

		requisition.SMSReplyCount++
		if _, err = s.repository.UpdateRequisitionSMSReplyCount(requisition); err != nil {
			log.Printf("(ERR) Failed to update requisition [%s] sms reply count: %v", err)
		}
	}

	return nil
}
