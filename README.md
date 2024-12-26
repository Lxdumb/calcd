>## **Calcd - веб-сервис для вычисления математических выражений на Go**
__Calcd__ поддерживает следующие математические операции:  
* сложение (+)
* вычитание (-)
* умножение (*)
* деление (/)
* возведение в степень (^)
* скобки для действий  

__Calcd__ работает на порте 8080, т.е. URL веб-сервиса: `http://localhost:8080/api/v1/calculate`  
Пример использования (используется curl, встроена в большинство Linux-дистрибутивов, есть [версия под Windows](https://curl.se/windows/)):  
> `$ go run ./cmd/main.go`

В другом окне терминала:   
> `$ curl --location 'http://localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{"expression": "(3.1415926 * 8) / (6 + 3)"}'`  
> `{`  
> `"result":"2.7925267555555555"`  
> `}`

Примеры ошибок:  
>Ошибка 422:  
> > `$ curl --location 'http://localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{"expression": "(abc)"}'`  
> `{`  
> `"error": "Expression is not valid"`  
> `}`  

>Ошибка 500:  
> > `$ curl --location 'http://localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{"expression": "\(*_*)/"}'`  
> `{`  
> `"error": "Internal server error"`  
> `}`

### Как это работает?  
Для вычисления математической формулы используется метод преобразования выражения в RPN (обратная польская запись), где не используются скобки, а затем уже вычислить новое выражение. Подробнее о RPN [тут](https://ru.wikipedia.org/wiki/%D0%9E%D0%B1%D1%80%D0%B0%D1%82%D0%BD%D0%B0%D1%8F_%D0%BF%D0%BE%D0%BB%D1%8C%D1%81%D0%BA%D0%B0%D1%8F_%D0%B7%D0%B0%D0%BF%D0%B8%D1%81%D1%8C)  
Для приема запросов и их отправки используются встроенные в Go функции и библиотеки.
