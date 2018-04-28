mongo-sync-go
Inspired from https://github.com/sheharyarn/mongo-sync

***Setup***
1. Create a config file `~/.mongo-sync.yaml` by copying provided template file: `.mongo-sync.yaml.template`
2. Update `~/.mongo-sync.yaml` file to to provide remote and local ( could be another remote database) mongo server configuration

***Installing mongo-sync-go***
1. Using homebrew
    ```
    brew install bharat-p/tap/mongo-sync-go
    ```

2. Download pre-built binary from:           

    https://github.com/bharat-p/mongo-sync-go/releases


3. build your own binary
    ```
   go get github.com/bharat-p/mongo-sync-go

   cd $GOPATH/src/github.com/bharat-p/mongo-sync-go

   go build -o mongo-sync-go main.go 
   ```

***Usage***
1.  ```
    mongo-sync-go pull --remote.database=mydb
    ```

    Sync/copy database: `mydb` from remote server into database name `mydb` in local server

    _Remote and local database credentials/settings will be used from config file: `~/.mongo-sync.yaml`_
    
2. ```
    mongo-sync-go pull --remote.database=mydb --local.database=my-local-db
    ```

   Sync/copy database: `mydb` from remote server into database name `my-local-db` in local server

3. ```
   mongo-sync-go pull --remote.host=mongo.remote --local.host=mongo.local --remote.database=mydb
   ```

   Sync/copy database: `mydb` from host: `mongo.remote` to database name `my-db` onto server `mongo.local`

4. ```
   mongo-sync-go pull -h
   ``` 
   to see all available options

