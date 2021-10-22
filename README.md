# golang-roles-permissions
### roles&permissions
- roles have permissions
- users have roles
- users can have granted or revoked permissions
### general info
- repository pattern
- custom validation package (dtos, query strings)
- stateful tokens (fast hashed with sha256)
- two types of json responses ok and error 
- pagination with metadata
- rate limiting
- graceful shutdown

### middlewares
- is user anonymous/authenticated?
- is user activated?
- does user have permission?

### tools
- net/http
- gorm
- zap jsonlogger (sugar)
