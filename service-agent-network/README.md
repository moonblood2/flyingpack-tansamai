# Service Agent Network
A REST API service for agent network.

## Technical Terms
### Stacks that use in this project includes
* Typescript for main language.
* Express framework for handle REST API request and response.
* Postgresql for database management and there is only the raw SQL, not ORM.

### Design pattern and project layout.
* controller:
* service:
* data access object (dao):
* data transfer object (dto):
  You can read it further here https://www.toptal.com/express-js/nodejs-typescript-rest-api-pt-1

### Deployment
Both staging and production use Heroku dyno to run the server and use Heroku Postgres for the database.

## Dev
### Run database
```docker-compose -f docker-compose.yml up -d postgres```
### Run server
* Run with "nodemon" it will rebuild and start node every time when source code in the project change.
  ```npm run dev```
## Build and run
* You can use tsc and start with node.
  build: ```tsc ./src/app.ts```
  run: ```node ./dist/app.js```
