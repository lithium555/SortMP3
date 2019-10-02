# SortMP3

#Up postgre in container
```
docker run --name googleAPI_new -e POSTGRES_PASSWORD=rem -p 5432:5432 -d postgres

docker run --name sort_music -e POSTGRES_PASSWORD=master -e POSTGRES_DB=musicDB -e POSTGRES_USER=sorter -p 5432:5432 -d postgres

53aa124bb2af  postgres    "docker-entrypoint.s…"  About a minute ago   Up About a minute   0.0.0.0:5432->5432/tcp   googleAPI_new

docker rm -f sort_music

```


 >>>>>>>> лучше сначала продумай как объекты вообще связаны забывая что там есть папки
>>>>>>>> то есть какие у них связи в реальном мире

>>>>>>>>в идеале опиши это как CREATE TABLE, но пока без FOREIGN KEY
>>>>>>>>так будет точно понятно что ты имеешь в виду
Я перенес всю переписку (там не много) в текстовый файл.
Следующий шаг  это ответиь на твой вопрос :
 >>>>>>>> лучше сначала продумай как объекты вообще связаны забывая что там есть папки
>>>>>>>> то есть какие у них связи в реальном мире

Как я это понял: в реальном мире конечная таблица должна быть такая :
 genre | Name of Singer | album Name | Name of composition
Чтобы ее получить надо понять кто по иерархии от кого зависит :
Песня зависит от альбома,
альбом зависит от автора,  песня зависит от жанра.
Так же песня завиист от автора
Получаем такую типа диаграмму:
Song -> Album
Song -> Author 
Song -> Genre

Album -> Author
Это можно представить, как 4 маленьких таблички:

Жанр:
GENRE
------------------ 
id
name


Автор:

AUTHOR
_______________
id
author_name


Альбом:
 
ALBUM
_______________
id
authorID
album_name
year (дата выхода альбома), 
cover (ссылка или путь к картинке альбома)



Песня:
 
SONG
-------------------
id
name
albumID
genreID
authorID
trackNum (номер песни в альбоме)


Так же я понял как пофиксить регистр, когда оодна и та ж епесня может попадать с большой и маленькой буквы в названии.
Надо все имена песен исполнителей и так далее опускать в нижний регистр. Так мы обеспечим уникальность для каждой песни и не будет повторов.
А вот что дальше я пока не знаю.
Пока погуглю как создать пакет с postgreSQL и добавлю его в проект.


ID3 теги: 
https://en.wikipedia.org/wiki/ID3
либы:
https://github.com/mikkyang/id3-go
или
https://github.com/dhowden/tag


Про почему мало инфы - можно поставить тот же Picard и сравнить с тем что он видит. 
Может просто в файле мало всего записано.
https://picard.musicbrainz.org/

https://github.com/mikkyang/id3-go
https://github.com/dhowden/tag


https://stackoverflow.com/questions/24712463/go-there-is-no-parameter-1
