# projektowanie-echo-go


Aby uruchomić program główny - server.go - należy uruchomić w terminalu w jego katalogu nadrzędnym:
-> go run server.go 


Wówczas na porcie 8000 dostępne jest API do dodawania, usuwania, aktualizowania oraz pobierania obiektów "items" z bazy myItems.db
Obiekty składają się z ID oraz z Nazwy Name. Przy tworzeniu nowego obiektu "item" należy podać jedynie nazwę.
Testowanie odbywa się z pomocą programu Postman
Przykładowe body do createItem (POST) oraz updateItem (PUT)

POST: 

 {
 
  "name": "Name_sample"
  
 }
 
http://localhost:11223/items

Odpowiedź:

{

    "id": 3,
    "name": "Name_sample"
    
}


 
GET: 

GET http://localhost:11223/items/1 

 Odpowiedź 
 
{

    "id": 1,
    "name": "Name_sample"
    
}

PUT: 

body: {

    "name": "New_name"
    
}

PUT http://localhost:11223/items/2


Odpowiedź 

{

    "id": 2,
    "name": "New_name"
    
}

DELETE: 

DELETE http://localhost:11223/items/1


Po wykonaniu operacji item o danym ID nie jest już dostępny
