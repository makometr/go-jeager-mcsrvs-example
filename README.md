# go-jeager-mcsrvs-example

## Сервис Worker

Имеет одну ручку скдадывания чисел. Генерирует ошибки, если: \

- не смог распарсить json
- передан массив чисел длиной > 5
- в массиве есть нули

Особенности: \

- каждое число "складывается" 100мс
- каждая единица "складывается" 3с
