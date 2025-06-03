package cron

import (
	"github.com/robfig/cron/v3"
	"log"
)

type handler struct {
	usecase Usecase
}

func Init(usecase Usecase) {
	h := handler{usecase: usecase}
	c := cron.New()

	if _, err := c.AddFunc("@daily", func() {
		if err := h.usecase.eventUpdateRemarkSoonToOpen(); err != nil {
			log.Println("update remark event soon to open cron : " + err.Error())
		} else {
			log.Println("running task remark event update soon to open success")
		}
	}); err != nil {
		log.Fatal(err)
	}

	if _, err := c.AddFunc("@daily", func() {
		if err := h.usecase.eventUpdateRemarkOpenToClose(); err != nil {
			log.Println("update remark event open to close cron : " + err.Error())
		} else {
			log.Println("running task remark event update open to close success")
		}
	}); err != nil {
		log.Fatal(err)
	}
	if _, err := c.AddFunc("@daily", func() {
		if err := h.usecase.eventUpdateRemarkCloseToOngoing(); err != nil {
			log.Println("update remark event close to ongoing cron : " + err.Error())
		} else {
			log.Println("running task remark event update close to ongoing success")
		}
	}); err != nil {
		log.Fatal(err)
	}
	//c.AddFunc("@daily", func() {
	//	if err := h.usecase.eventUpdateRemarkOngoingToDone(); err != nil {
	//		log.Println("update remark event ongoing to done cron : " + err.Error())
	//	} else {
	//		log.Println("running task remark event update ongoing to done success")
	//	}
	//})
	if _, err := c.AddFunc("@every 1h", func() {
		if err := h.usecase.registrationUpdate(); err != nil {
			log.Println("update status registration user error : " + err.Error())
		} else {
			log.Println("running task status registration user success")
		}
	}); err != nil {
		log.Fatal(err)
	}
	c.Start()
}
