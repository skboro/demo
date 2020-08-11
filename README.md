# sellerapp assignment

Database used : MySQL
Build and run both services individually.
Verify the config files for the db connection info.

## User management and auth
- `demo-user-mgmt` is the service
- also serves as rudimentary reverse proxy for `demo-auction`
- supported APIs:
    1. /signup : signup a new user
    ```
        curl 'http://localhost:3001/signup' -H 'content-type: application/json' --data-raw '{"name":"suman","email":"suman@a.com","password":"password"}'
    ```
    2. /signin : signin as a user, returns JWT token to be used for latter APIs. Also sets the token as cookie.
    ```
        curl 'http://localhost:3001/signin' -H 'content-type: application/json' --data-raw '{"email":"suman@a.com","password":"password"}'

        # admin can login with predefined secret key
        curl 'http://localhost:3001/signin' -H 'content-type: application/json' --data-raw '{"email":"admin@sellerapp.com","password":"admin_secret_key"}'
    ```
    3. /update : update account for current user
    4. /delete : delete account for current user
    5. /getAllUsers : admin api to get list of all users
    ```
        curl 'http://localhost:3001/getAllUsers' -H 'cookie: token="eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowLCJleHAiOjE1OTcyMjUzMDl9.r_3VhoJ35_2UJM02GjJlqfk_J5sq0mV4V1iR8g9ajwIRSzdJ5vSSv2PTVzb6-kB-hsvYXOcuHA5kFdGALto7Cw"'
    ```
    6. /account : returns the current user account info
    ```
        curl 'http://localhost:3001/account' -H 'cookie: token="eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE1OTcyMjUyNTN9.Qmh7GjgFvRdgHAd_4NeMETuszRJxAEis-PikjZKBSGiQh4vOKQMzcnxGSEZTjy-dwPa-2ss_-noXJwwgg87LaA"'
    ```

Sample `.config` file
```
{
    "port": 3001,
    "auction_url": "localhost:3002",
    "jwt_key": "secret_key",
    "database": {
        "host": "localhost",
        "username": "user",
        "password": "password",
        "dbname": "mydb"
    }
}
```
## Auction management and bidding
- `demo-auction` is the service
- supported APIs:
    1. /auction/create : create an auction, admin api
    ```
        curl http://localhost:3001/auction/create -H 'content-type: application/json' -H 'cookie: token="eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowLCJleHAiOjE1OTcyMjUzMDl9.r_3VhoJ35_2UJM02GjJlqfk_J5sq0mV4V1iR8g9ajwIRSzdJ5vSSv2PTVzb6-kB-hsvYXOcuHA5kFdGALto7Cw"' --data-raw '
        {
            "start_time": "2020-08-11T15:47:00.000000+05:30",
            "end_time": "2020-08-11T15:49:00.000000+05:30",
            "start_price": 1200,
            "item_name": "t-shirt"
        }'
    ```
    2. /auction/update : update an auction, admin api
    3. /auction/delete : delete an auction, admin api
    4. /auction/getAll : get all auctions, admin api
    ```
        curl http://localhost:3001/auction/getAll -H 'cookie: token="eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowLCJleHAiOjE1OTcyMjUzMDl9.r_3VhoJ35_2UJM02GjJlqfk_J5sq0mV4V1iR8g9ajwIRSzdJ5vSSv2PTVzb6-kB-hsvYXOcuHA5kFdGALto7Cw"'
    ```
    5. /auction/getLive : get live auctions
    ```
        curl http://localhost:3001/auction/getLive -H 'cookie: token="eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE1OTcyMjUyNTN9.Qmh7GjgFvRdgHAd_4NeMETuszRJxAEis-PikjZKBSGiQh4vOKQMzcnxGSEZTjy-dwPa-2ss_-noXJwwgg87LaA"'
    ```
    6. /auction/bid/create : submit a bid
    ```
        curl 'http://localhost:3001/auction/bid/create' -H 'cookie: token="eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE1OTcyMjUyNTN9.Qmh7GjgFvRdgHAd_4NeMETuszRJxAEis-PikjZKBSGiQh4vOKQMzcnxGSEZTjy-dwPa-2ss_-noXJwwgg87LaA"' -H 'content-type: application/json' --data-raw '
        {
            "price": 1220,
            "user_id": 1,
            "auction_id": 2
        }'
    ```
    7. /auction/bid/update : update a bid
    8. /auction/bid/delete : delete a bid
    9. /auction/bid/get : get bid info
    ```
        curl 'http://localhost:3001/auction/bid/get' -H 'cookie: token="eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE1OTcyMjUyNTN9.Qmh7GjgFvRdgHAd_4NeMETuszRJxAEis-PikjZKBSGiQh4vOKQMzcnxGSEZTjy-dwPa-2ss_-noXJwwgg87LaA"' -H 'content-type: application/json' --data-raw '{"id": 1}'
    ```
    10. /auction/bid/getBids :
    ```
    # get all bids, admin api
        curl 'http://localhost:3001/auction/bid/getBids' -H 'cookie: token="eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowLCJleHAiOjE1OTcyMjUzMDl9.r_3VhoJ35_2UJM02GjJlqfk_J5sq0mV4V1iR8g9ajwIRSzdJ5vSSv2PTVzb6-kB-hsvYXOcuHA5kFdGALto7Cw"' -H 'content-type: application/json' --data-raw '{}'
    
    # get all bids of an auction, admin api
        curl 'http://localhost:3001/auction/bid/getBids' -H 'cookie: token="eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowLCJleHAiOjE1OTcyMjUzMDl9.r_3VhoJ35_2UJM02GjJlqfk_J5sq0mV4V1iR8g9ajwIRSzdJ5vSSv2PTVzb6-kB-hsvYXOcuHA5kFdGALto7Cw"' -H 'content-type: application/json' --data-raw '{"auction_id":1}'
    
    # all bids by a user
        curl 'http://localhost:3001/auction/bid/getBids' -H 'cookie: token="eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE1OTcyMjUyNTN9.Qmh7GjgFvRdgHAd_4NeMETuszRJxAEis-PikjZKBSGiQh4vOKQMzcnxGSEZTjy-dwPa-2ss_-noXJwwgg87LaA"' -H 'content-type: application/json' --data-raw '{"user_id":1}'
    ```

Sample `.auction_config` file
```
{
    "port": 3002,
    "database": {
        "host": "localhost",
        "username": "user",
        "password": "password",
        "dbname": "mydb"
    }
}
```