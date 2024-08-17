# AUTHORS:

Andreas Lootus / Andreas.Lootus  
Iisak Virma / ivirma  
Johannes Sild / jsild

Project git url: https://01.kood.tech/git/jsild/social-network  
URL to Audit page: https://github.com/01-edu/public/blob/master/subjects/social-network/audit  

## INSTRUCTIONS
In order to make things easier for you, we have made a bash file. It runs the necessary commands to build and run the docker images.

- Docker
```bash
chmod 777 ./runDocker.sh
./runDocker.sh
```
In your web browser, navigate to the following URL: http://localhost:5173/
  
OR you can do it manually by:  
  
- Manual  
Install all nessecary packages with 
```bash
npm install
```

Open 2 terminals  
In terminal 1 run:  
```bash
cd backend/server/ && go run .
```

In terminal 2 run:  
```bash
cd vue-frontend && npm run dev
```
In your web browser, navigate to the following URL: http://localhost:5173/

## DESCRIPTION:

Social-network is a facebook-like webproject.
As a base we used real-time-forum (backend mostly) and chose to use vue.js for frontend's framework.  
We have seperated back and frontend servers, both run on a different server and communicate through websockets. For migrations we used golang-migrations, migrations are located in backend/internal/data/migrations/sqlite3 and are ran when the backend server is started.  
