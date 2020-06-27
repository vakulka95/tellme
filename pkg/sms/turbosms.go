package sms

import (
	"fmt"
	"log"
	"time"

	"gitlab.com/tellmecomua/tellme.api/app/representation"

	turbosms "github.com/wildsurfer/turbosms-go"
)

const (
	RequisitionApplyTemplate = `Вашу заявку успішно зареєстровано на сайті "Розкажи мені". 

Зараз сайт працює в тестовому режимі, тому час консультації може змінитися.  Про всі зміни часу консультації ви можете домовитися з психологом.`

	RequisitionReviewTemplate = `Оцініть, будь ласка, Ваш досвід користування платформою «Розкажи мені»: %s`

	RequisitionReplyMaleTemplate   = `Шановний %s, ми поставили Вашу заявку у список очікування. Але не турбуйтесь, уже найближчим часом ми візьмемо її в роботу і з Вами обов’язково зв'яжеться психолог.`
	RequisitionReplyFemaleTemplate = `Шановна %s, ми поставили Вашу заявку у список очікування. Але не турбуйтесь, уже найближчим часом ми візьмемо її в роботу і з Вами обов’язково зв'яжеться психолог.`
)

var genderToRequisitionReplyTemplate = map[string]string{
	representation.GenderMale:   RequisitionReplyMaleTemplate,
	representation.GenderFemale: RequisitionReplyFemaleTemplate,
}

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

func (s *TurboSMS) SendRequisitionReply(phone, username, gender string) error {
	resp, err := s.cli.SendSMS("Tell me", "+38"+phone, fmt.Sprintf(genderToRequisitionReplyTemplate[gender], username), "")
	if err != nil {
		return err
	}

	log.Printf("(INFO) Send requisition reply to [%s] done: %+v", phone, resp.SendSMSResult)
	return nil
}
