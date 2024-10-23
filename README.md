# Обновление
1. Вынес логику работы с БД из [handlers.go](internal/handlers/handlers.go#L30) в [employeesRepository.go](internal/repositories/employeesRepository/employeesRepository.go)
2. Разобрался с тем что такое middleware, переписал авторизацию [authorization.go](internal/middleware/authorization.go)
3. Сделал динамическое построение запроса для обновления [updateQuery.go](internal/dbservice/updateQuery.go#L43)
4. Закешировал подключение к БД в [employeesRepository.go](internal/repositories/employeesRepository/employeesRepository.go#L26) в него же перенес валидацию 
5. Исправил остальные недочеты: 
    
    5.1 проверка [ошибки](internal/dbservice/dbservice.go#L149) после итерации по строкам
    
    5.2 вынес [конфиг](internal/dbservice/dbservice.go#L15) бд в переменные среды
    
    5.3 убрал wildcard из [SELECT](internal/dbservice/dbservice.go#L107)
