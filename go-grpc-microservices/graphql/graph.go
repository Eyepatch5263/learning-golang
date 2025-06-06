package main

// type Server struct {
// 	accountClient	*account.client
// 	productClient	*product.client
// 	orderClient		*order.client
// }

func NewGraphQLServer(accountUrl,productUrl,orderUrl string) (*Server,error) {
	accountClient,err:=account.NewClient(accountUrl)
	if err!=nil {
		return nil,err
	}

	productClient,err:=account.NewClient(productUrl)
	if err!=nil {
		accountClient.close()
		return nil,err
	}

	orderClient,err:=account.NewClient(orderUrl)
	if err!=nil {
		accountClient.close()
		productClient.close()
		return nil,err
	}

	


}