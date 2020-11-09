# Readme
REST API
* `/login` - получаем токены по json
* `/reg` - Страница для регистрации юзера
* `/refresh` - Выполняет операцию Refresh, вы получаете новые токены Access и Refresh
* `/logout` - Удаляет текущий Refresh токен
* `/logoutAll` - Удаляет все Токены юзера

Сервер и клиент работают по json, используется JWT для Refresh и Access токенов.
Для получения токена, клиент отправляет GUID и данные юзера ` {"name": username, "password" : password, "uuid" :  GUID}`,
 после успешной аутентификации клиенту передается Access и Refresh токен.

