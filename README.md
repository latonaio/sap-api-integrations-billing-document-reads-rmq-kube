# sap-api-integrations-billing-document-reads-rmq-kube 
sap-api-integrations-billing-document-reads-rmq-kube は、外部システム(特にエッジコンピューティング環境)をSAPと統合することを目的に、SAP API で請求伝票データ を取得するマイクロサービスです。    
sap-api-integrations-billing-document-reads-rmq-kube には、サンプルのAPI Json フォーマットが含まれています。   
sap-api-integrations-billing-document-reads-rmq-kube は、オンプレミス版である（＝クラウド版ではない）SAPS4HANA API の利用を前提としています。クラウド版APIを利用する場合は、ご注意ください。   
https://api.sap.com/api/OP_API_BILLING_DOCUMENT_SRV_0001/overview

## 動作環境  
sap-api-integrations-billing-document-reads-rmq-kube は、主にエッジコンピューティング環境における動作にフォーカスしています。  
使用する際は、事前に下記の通り エッジコンピューティングの動作環境（推奨/必須）を用意してください。  
・ エッジ Kubernetes （推奨）    
・ AION のリソース （推奨)    
・ OS: LinuxOS （必須）    
・ CPU: ARM/AMD/Intel（いずれか必須）  
・ RabbitMQ on Kubernetes  
・ RabbitMQ Client      

## クラウド環境での利用
sap-api-integrations-billing-document-reads-rmq-kube は、外部システムがクラウド環境である場合にSAPと統合するときにおいても、利用可能なように設計されています。  

## 本レポジトリ が 対応する API サービス
sap-api-integrations-billing-document-reads-rmq-kube が対応する APIサービス は、次のものです。

* APIサービス概要説明 URL: https://api.sap.com/api/OP_API_BILLING_DOCUMENT_SRV_0001/overview  
* APIサービス名(=baseURL): API_BILLING_DOCUMENT_SRV

## RabbitMQ からの JSON Input

sap-api-integrations-billing-document-reads-rmq-kube は、Inputとして、RabbitMQ からのメッセージをJSON形式で受け取ります。 
Input の サンプルJSON は、Inputs フォルダ内にあります。  

## RabbitMQ からのメッセージ受信による イベントドリヴン の ランタイム実行

sap-api-integrations-billing-document-reads-rmq-kube は、RabbitMQ からのメッセージを受け取ると、イベントドリヴンでランタイムを実行します。  
AION の仕様では、Kubernetes 上 の 当該マイクロサービスPod は 立ち上がったまま待機状態で当該メッセージを受け取り、（コンテナ起動などの段取時間をカットして）即座にランタイムを実行します。　

## RabbitMQ への JSON Output

sap-api-integrations-billing-document-reads-rmq-kube は、Outputとして、RabbitMQ へのメッセージをJSON形式で出力します。  
Output の サンプルJSON は、Outputs フォルダ内にあります。  

## RabbitMQ の マスタサーバ環境

sap-api-integrations-billing-document-reads-rmq-kube が利用する RabbitMQ のマスタサーバ環境は、[rabbitmq-on-kubernetes](https://github.com/latonaio/rabbitmq-on-kubernetes) です。  
当該マスタサーバ環境は、同じエッジコンピューティングデバイスに配置されても、別の物理(仮想)サーバ内に配置されても、どちらでも構いません。

## RabbitMQ の Golang Runtime ライブラリ
sap-api-integrations-billing-document-reads-rmq-kube は、RabbitMQ の Golang Runtime ライブラリ として、[rabbitmq-golang-client](https://github.com/latonaio/rabbitmq-golang-client)を利用しています。

## デプロイ・稼働
sap-api-integrations-billing-document-reads-rmq-kube の デプロイ・稼働 を行うためには、aion-service-definitions の services.yml に、本レポジトリの services.yml を設定する必要があります。

kubectl apply - f 等で Deployment作成後、以下のコマンドで Pod が正しく生成されていることを確認してください。
```
$ kubectl get pods
```

## 本レポジトリ に 含まれる API名
sap-api-integrations-billing-document-reads-rmq-kube には、次の API をコールするためのリソースが含まれています。  

* A_BillingDocument（請求伝票 - ヘッダ）※請求伝票の詳細データを取得するために、ToItem、ToHeaderPartner、ToItemPartner、ToItemPricingElement、と合わせて利用されます。
* A_BillingDocumentPartner（請求伝票 - ヘッダ取引先）
* A_BillingDocumentItem（請求伝票 - 明細）※請求伝票明細の詳細データを取得するために、ToItemPartner、ToItemPricingElementと合わせて利用されます。
* A_BillingDocumentItemPartner（請求伝票 - 明細取引先）
* ToHeaderPartner（請求伝票 - ヘッダ取引先 ※To）
* ToItem（請求伝票 - 明細 ※To）
* ToItemPartner（請求伝票 - 明細取引先 ※To）
* ToItemPricingElement（請求伝票 - 明細価格条件）

## API への 値入力条件 の 初期値
sap-api-integrations-billing-document-reads-rmq-kube において、API への値入力条件の初期値は、入力ファイルレイアウトの種別毎に、次の通りとなっています。  

### SDC レイアウト

* inoutSDC.BillingDocument.BillingDocument（請求伝票）
* inoutSDC.BillingDocument.HeaderPartner.PartnerFunction（ヘッダ取引先機能）
* inoutSDC.BillingDocument.BillingDocumentItem.BillingDocumentItem（請求伝票明細）
* inoutSDC.BillingDocument.BillingDocumentItem.ItemPartner.PartnerFunction（明細取引先機能）

## SAP API Bussiness Hub の API の選択的コール

Latona および AION の SAP 関連リソースでは、Inputs フォルダ下の sample.json の accepter に取得したいデータの種別（＝APIの種別）を入力し、指定することができます。  
なお、同 accepter にAll(もしくは空白)の値を入力することで、全データ（＝全APIの種別）をまとめて取得することができます。  

* sample.jsonの記載例(1)  

accepter において 下記の例のように、データの種別（＝APIの種別）を指定します。  
ここでは、"Header"が指定されています。    
  
```
	"api_schema": "sap.s4.beh.billingdocument.v1.BillingDocument.Created.v1",
	"accepter": ["Header"],	
	"billing_document": "90000000",
	"deleted": false
```
  
* 全データを取得する際のsample.jsonの記載例(2)  

全データを取得する場合、sample.json は以下のように記載します。  

```
	"api_schema": "sap.s4.beh.billingdocument.v1.BillingDocument.Created.v1",
	"accepter": ["All"],	
	"billing_document": "90000000",
	"deleted": false
```

## 指定されたデータ種別のコール

accepter における データ種別 の指定に基づいて SAP_API_Caller 内の caller.go で API がコールされます。  
caller.go の func() 毎 の 以下の箇所が、指定された API をコールするソースコードです。  

```
func (c *SAPAPICaller) AsyncGetBillingDocument(billingDocument, headerPartnerFunction, billingDocumentItem, itemPartnerFunction string, accepter []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(accepter))
	for _, fn := range accepter {
		switch fn {
		case "Header":
			func() {
				c.Header(billingDocument)
				wg.Done()
			}()
		case "HeaderPartner":
			func() {
				c.HeaderPartner(billingDocument, headerPartnerFunction)
				wg.Done()
			}()
		case "Item":
			func() {
				c.Item(billingDocument, billingDocumentItem)
				wg.Done()
			}()
		case "ItemPartner":
			func() {
				c.ItemPartner(billingDocument, billingDocumentItem, itemPartnerFunction)
				wg.Done()
			}()
		default:
			wg.Done()
		}
	}
	wg.Wait()
}
```

## Output  
本マイクロサービスでは、[golang-logging-library-for-sap](https://github.com/latonaio/golang-logging-library-for-sap) により、以下のようなデータがJSON形式で出力されます。  
以下の sample.json の例は、SAP 請求伝票 の ヘッダデータ が取得された結果の JSON の例です。  
以下の項目のうち、"BillingDocument" ～ "to_Partner" は、/SAP_API_Output_Formatter/type.go 内 の Type Header {} による出力結果です。"cursor" ～ "time"は、golang-logging-library による 定型フォーマットの出力結果です。  

```
{
	"cursor": "/Users/latona2/bitbucket/sap-api-integrations-billing-document-reads/SAP_API_Caller/caller.go#L67",
	"function": "sap-api-integrations-billing-document-reads/SAP_API_Caller.(*SAPAPICaller).Header",
	"level": "INFO",
	"message": [
		{
			"BillingDocument": "90000000",
			"BillingDocumentType": "F2",
			"SDDocumentCategory": "M",
			"BillingDocumentCategory": "L",
			"CreationDate": "2016-08-17T09:00:00+09:00",
			"LastChangeDate": "",
			"SalesOrganization": "1710",
			"DistributionChannel": "10",
			"Division": "00",
			"BillingDocumentDate": "2016-08-17T09:00:00+09:00",
			"BillingDocumentIsCancelled": false,
			"CancelledBillingDocument": "",
			"IsExportDelivery": "",
			"TotalNetAmount": "52.65",
			"TransactionCurrency": "USD",
			"TaxAmount": "4.21",
			"TotalGrossAmount": "56.86",
			"CustomerPriceGroup": "",
			"IncotermsClassification": "EXW",
			"CustomerPaymentTerms": "0004",
			"PaymentMethod": "",
			"CompanyCode": "1710",
			"AccountingDocument": "",
			"ExchangeRateDate": "2016-08-17T09:00:00+09:00",
			"ExchangeRateType": "",
			"DocumentReferenceID": "0090000000",
			"SoldToParty": "17100001",
			"PartnerCompany": "",
			"PurchaseOrderByCustomer": "",
			"CustomerGroup": "01",
			"Country": "US",
			"CityCode": "",
			"Region": "GA",
			"CreditControlArea": "A000",
			"OverallBillingStatus": "A",
			"AccountingPostingStatus": "C",
			"AccountingTransferStatus": "C",
			"InvoiceListStatus": "",
			"BillingDocumentListType": "LR",
			"BillingDocumentListDate": "",
			"to_Item": "https://sandbox.api.sap.com/s4hanacloud/sap/opu/odata/sap/API_BILLING_DOCUMENT_SRV/A_BillingDocument('90000000')/to_Item",
			"to_Partner": "https://sandbox.api.sap.com/s4hanacloud/sap/opu/odata/sap/API_BILLING_DOCUMENT_SRV/A_BillingDocument('90000000')/to_Partner"
		}
	],
	"time": "2022-01-27T23:01:13+09:00"
}
```