# taskman backend

## TODO

### backend

#### general

...


#### user

- [x] account email field (case insensitive)
- [x] get user by id
- [x] edit user (username & password)
    - [ ] reset password
- [ ] delete user
- [x] login (return jwt)


#### tasks 

**task**
- description
- due date
- completion (todo, in_progress, done)

- [x] new table  
    - each user gets a table
    - how to authenticate?
    - **Login returns a JWT in the JSON response**
        - in requests to protected endpoints (middleware), include the JWT token in the `Authorization` header
        - `Authorization: Bearer <JWT token>`
    - holds tasks
    - foreign key account_id that references accounts(id)

- [ ] new task
- [ ] list tasks
- [ ] delete task
- [ ] edit task


### misc

- [ ] refactor from Gin to net/http?
- [ ] Github OAuth2.0
    - https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/creating-an-oauth-app
