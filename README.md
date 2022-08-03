## v1.0.2

#### Introduction to RedisCache
RedisCache is my personal Go project, which enables multiple clients to get, set, delete and update data

#### Setup
1. Install Docker and Docker Desktop on your computer: [Installation link](https://www.google.com/url?sa=t&rct=j&q=&esrc=s&source=web&cd=&cad=rja&uact=8&ved=2ahUKEwjdv6CE6qr5AhUFD-wKHY00AdUQFnoECAYQAQ&url=https%3A%2F%2Fdocs.docker.com%2Fengine%2Finstall%2F&usg=AOvVaw3oxUtu6GW_HNWz3ZCPMLU_)
2. Install Redis Extension in Docker Desktop
3. Start redis by running `docker run --name my-redis -p 6379:6379 -d redis` on your terminal

#### Get Started 
1. Clone the repository: `git clone https://github.com/yamaceay/rediscache`
2. Change the current directory: `cd rediscache`
3. [Optional] Configure settings.json using "config" mode (-> Parameters) if you want to use another IP address
4. Make sure the Redis container is running on Docker Desktop (-> Containers)
5. Start the server: `go run .`
6. Starting a new client: 
  1. Open a new terminal in the same directory
  2. Run `go run . --mode=client` to see the commands

#### Settings 
If you want to change the behaviour of server, you can use the command:
`go run . -M config --dbAddress=<dbHost>:<dbPort> --ipAddress=<ipHost>:<ipPort> --ttlMinutes:<ttlMinutes>`

It simply creates a settings.json with the following environment variables: 

| Name | Description                     | Type | Options | Default |
| :----: | :-----------------------------: | :-: | :-------: | :-------: |
| dbHost | Database host | string | any ip address | "localhost" |
| dbPort | Database port | int | any ip port | 6379 |
| ipHost | Application host | string | any ip address | "localhost" |
| ipPort | Application port | int | any ip port | 8080 |
| ttlMinutes | Expiry time of data objects in minutes | int | any positive integer | 10080 |    

#### Arguments
These are the arguments that a client may use when sending a request using CLI program. 

| Name | Alias | Description                     | Type | Options | Default |
| :----: | :---: | :-----------------------------: | :-: | :-------: | :-------: |
| mode | M | Server side or client side  | string | "client", "server", "config" | "server"  |
| method | X | What to do | string | "", "get", "set", "delete" | "" |
| db | - | Which database | int | 0, 1, … | 0 |
| key, value | k, v | Key-value pairs | string | - | "" |

Also, you can directly launch the client using the function in modes package
`func StartClient(settings Settings, filters ...Filter) error`
but filters need to be set accordingly.
