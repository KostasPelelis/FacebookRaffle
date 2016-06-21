# Facebook Raffle

A simple program writen in Go that uses the facebook api in order to run raffles against users that shared a specific post.

## Instalation

(Requires a go compiler installed)
* Clone the repository
```
git clone https://github.com/KostasPelelis/FacebookRaffle.git
```
* Build The Project(for example with golang)
```
go build
```
* Run the Executable (In Unix)
```
./shareraffle
```
## Execution Details

The program starts an HTTP webserver that runs on port 8080 and responds only in requests done in URL http://localhost/raffle?_params_  
_params_ are the options of the raffle which can be any of the following:

Name              | Type    | Explanation                                                                   | Required  | Default
----------------- | ------- | ----------------------------------------------------------------------------- | --------  | -------
user_id           | String  | The facebook ID of the page. This can be found [here](http://findmyfbid.com/) | Yes       | -
post_id           | String  | The ID of the post. I will upload a screencast on how to find it              | Yes       | -
winners           | Integer | The number of the winners of the raffle                                       | Yes       | -
liked_users_only  | Boolean | True if you want the winners to have liked the page                           | No        | False
unique_winners    | Boolean | True if you want only unique winners                                          | No        | False

