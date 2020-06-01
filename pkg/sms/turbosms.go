package sms

import (
	"fmt"
	"log"
	"time"

	turbosms "github.com/wildsurfer/turbosms-go"
)

const (
	RequisitionApplyTemplate = `Вашу заявку успішно зареєстровано на сайті "Розкажи мені". 

Зараз сайт працює в тестовому режимі, тому час консультації може змінитися.  Про всі зміни часу консультації ви можете домовитися з психологом.`

	RequisitionReviewTemplate = `Дякуємо, що обрали платформу психологічної підтримки "Розкажи мені". Нам дуже важлива ваша думка про досвід користування платформою та спілкування з психологом. Просимо оцінити нашу роботу за посиланням: %s`
)

type TurboSMS struct {
	cli *turbosms.Client
}

func NewManager(username, password string) *TurboSMS {
	srv := &TurboSMS{
		cli: turbosms.NewClient(username, password),
	}

	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for {
			select {
			case <-ticker.C:
				srv.cli = turbosms.NewClient(username, password)
				log.Printf("(INFO) TurboSMS client reinitialized ")
			}
		}
	}()

	return srv
}

func (s *TurboSMS) SendRequisitionApply(phone string) error {
	resp, err := s.cli.SendSMS("Tell me", "+38"+phone, RequisitionApplyTemplate, "")
	if err != nil {
		return err
	}

	log.Printf("(INFO) Send requisition apply to [%s] done: %+v", phone, resp.SendSMSResult)
	return nil
}

func (s *TurboSMS) SendRequisitionReview(phone, link string) error {
	resp, err := s.cli.SendSMS("Tell me", "+38"+phone, fmt.Sprintf(RequisitionReviewTemplate, link), "")
	if err != nil {
		return err
	}

	log.Printf("(INFO) Send requisition review to [%s] done: %+v", phone, resp.SendSMSResult)
	return nil
}
