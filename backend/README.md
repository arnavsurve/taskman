# taskman backend

## TODO

### backend

#### user

- [x] account email field (case insensitive)
- [x] get user by id
- [x] edit user (username & password)
    - [ ] reset password
- [ ] delete user


#### tasks 

**task**
- description
- due date
- completion (todo, in_progress, done)

- [ ] new table  
    - [ ] each user gets a table
        - how to authenticate?
        - [ ] Github OAuth2.0
        - https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/creating-an-oauth-app
        - **currently AddUser returns a JWT in the JSON response**
    - holds tasks
    - foreign key account_id that references accounts(id)

- [ ] new task
- [ ] list tasks
- [ ] delete task
- [ ] edit task


### misc

- [ ] refactor from Gin to net/http?
