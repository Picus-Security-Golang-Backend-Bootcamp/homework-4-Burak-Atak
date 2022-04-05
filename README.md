## Homework - 4

For listing all books.\
`http://localhost:8090/list`


***
For buying book by id.\
`http://localhost:8090/buy/{id}/{count}`


***
For deleting book by id.\
`http://localhost:8090/delete/{id}`

***
For searcing books by query key.\
`http://localhost:8090/search/{query}`

***
For creating a book.\
`http://localhost:8090/create`\
```
{
		"BookName":      string,
		"StockCode":     string,
		"Isbn":          string,
		"PageNumber":    int,
		"Price":         int,
		"StockQuantity": int,
		"Author":        string
	}
