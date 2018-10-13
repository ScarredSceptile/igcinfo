Develop an online service that will allow users to browse information about IGC files. IGC is an international file format for soaring track files that are used by paragliders and gliders. The program will not store anything in a persistent storage. Ie. no information will be stored on the server side on a disk or database. Instead, it will store submitted tracks in memory. Subsequent API calls will allow the user to browse and inspect stored IGC files.

For the development of the IGC processing, you will use an open source IGC library for Go: goigc

The system must be deployed on either Heroku or Google App Engine, and the Go source code must be available for inspection by the teaching staff (read-only access is sufficient).

App will be running on https://secret-hollows-32754.herokuapp.com/

How view the app:

All of the application will be under /igcinfo/

/igcinfo/api will get you the information about the app

/igcinfo/api/igc is where you can POST an igc and GET all the ids of the igcs in the app

/igcinfo/api/igc/<id> to get the igc of a given id
  
/igcinfo/api/igc/<id>/<field> to get the field of an igc with the given id.
  
Available fields are:

  pilot
  
  glider
  
  glider_id
  
  track_length
  
  H_date
