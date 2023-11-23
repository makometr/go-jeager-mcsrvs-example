# go-jeager-mcsrvs-example

![image](https://user-images.githubusercontent.com/21985069/222393530-b22fc76d-dafb-4662-9095-c98f52d4bdc8.png)
![image](https://user-images.githubusercontent.com/21985069/222394012-d36b2ad8-f115-408b-8683-9cdfe65509e7.png)

## Сервис Worker

Имеет одну ручку скдадывания чисел. Генерирует ошибки, если: \

- не смог распарсить json
- передан массив чисел длиной > 5
- в массиве есть нули

Особенности: \

- каждое число "складывается" 100мс
- каждая единица "складывается" 3с
