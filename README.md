# uas - user authentication service

Login to postgres using terminal, and create new database using command:
# CREATE DATABASE ___ ;
Using pattern from .env.dist, create new file called .env and replace examples with real credentials.

# Githooks
on unix systems, run:
```
cp .githooks/* .git/hooks/
chmod +x .git/hooks/*
```

#Creating session keys
For generating new keys, create new empty project, outside our folder, create generator.go and paste this:

```
func main(){
	hashKey := securecookie.GenerateRandomKey(32)
	hashKey64 := base64.StdEncoding.EncodeToString(hashKey)
	fmt.Println("SESSION_AUTH_KEY = " ,hashKey64)
	blockKey := securecookie.GenerateRandomKey(32)
	blockKey64 := base64.StdEncoding.EncodeToString(blockKey)
	fmt.Println("SESSION_ENCRYPTION_KEY = ",blockKey64)
}
```

Copy the results in .env file
#migrations
Use : $ migrate -source file://path/to/migrations -database postgres://user:password@localhost:5432/database up/down to  run migrations
