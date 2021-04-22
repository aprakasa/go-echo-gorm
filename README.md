# Go REST CRUD with Echo & Gorm

Quick and dirty CRUD REST API with Go, Echo and Gorm.


Create project
```
curl -X POST -H 'Content-type:application/json' -d '{"name":"golang rest", "description":"just another golang rest api"}' localhost:3333/projects
```

Get projects
```
curl localhost:3333/projects
```

Get single project
```
curl localhost:3333/projects/1
```

Update project
```
curl -X PATCH -H 'Content-type:application/json' -d '{"description":"rest api with golang"}' localhost:3333/projects/1
```

Delete project
```
curl -X DELETE localhost:3333/projects/1
```