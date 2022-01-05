package main

import (
	sap_api_caller "sap-api-integrations-billing-document-reads-rmq-kube/SAP_API_Caller"
	sap_api_input_reader "sap-api-integrations-billing-document-reads-rmq-kube/SAP_API_Input_Reader"
	"sap-api-integrations-billing-document-reads-rmq-kube/config"

	"github.com/latonaio/golang-logging-library/logger"
	rabbitmq "github.com/latonaio/rabbitmq-golang-client"
	"golang.org/x/xerrors"
)

func main() {
	l := logger.NewLogger()
	conf := config.NewConf()
	rmq, err := rabbitmq.NewRabbitmqClient(conf.RMQ.URL(), conf.RMQ.QueueFrom(), conf.RMQ.QueueTo())
	if err != nil {
		l.Fatal(err.Error())
	}
	defer rmq.Close()

	caller := sap_api_caller.NewSAPAPICaller(
		conf.SAP.BaseURL(),
		conf.RMQ.QueueTo(),
		rmq,
		l,
	)

	iter, err := rmq.Iterator()
	if err != nil {
		l.Fatal(err.Error())
	}
	defer rmq.Stop()

	for msg := range iter {
		err = callProcess(caller, msg)
		if err != nil {
			msg.Fail()
			l.Error(err)
			continue
		}
		msg.Success()
	}
}

func callProcess(caller *sap_api_caller.SAPAPICaller, msg rabbitmq.RabbitmqMessage) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = xerrors.Errorf("error occurred: %w", e)
			return
		}
	}()
	billingDocument, headerPartnerFunction, billingDocumentItem, itemPartnerFunction := extractData(msg.Data())
	accepter := getAccepter(msg.Data())
	caller.AsyncGetBillingDocument(billingDocument, headerPartnerFunction, billingDocumentItem, itemPartnerFunction, accepter)
	return nil
}

func extractData(data map[string]interface{}) (billingDocument, headerPartnerFunction, billingDocumentItem, itemPartnerFunction string) {
	sdc := sap_api_input_reader.ConvertToSDC(data)
	billingDocument = sdc.BillingDocument.BillingDocument
	headerPartnerFunction = sdc.BillingDocument.HeaderPartner.PartnerFunction
	billingDocumentItem = sdc.BillingDocument.BillingDocumentItem.BillingDocumentItem
	itemPartnerFunction = sdc.BillingDocument.BillingDocumentItem.ItemPartner.PartnerFunction
	return
}

func getAccepter(data map[string]interface{}) []string {
	sdc := sap_api_input_reader.ConvertToSDC(data)
	accepter := sdc.Accepter
	if len(sdc.Accepter) == 0 {
		accepter = []string{"All"}
	}

	if accepter[0] == "All" {
		accepter = []string{
			"Header", "HeaderPartner", "Item", "ItemPartner",
		}
	}
	return accepter
}
