### Introduction to RedisCache
RedisCache is my personal Go project, which enables multiple clients to get, set, delete and update data

### Setup
1. Install Docker and Docker Desktop on your computer: [Installation link](https://www.google.com/url?sa=t&rct=j&q=&esrc=s&source=web&cd=&cad=rja&uact=8&ved=2ahUKEwjdv6CE6qr5AhUFD-wKHY00AdUQFnoECAYQAQ&url=https%3A%2F%2Fdocs.docker.com%2Fengine%2Finstall%2F&usg=AOvVaw3oxUtu6GW_HNWz3ZCPMLU_)
2. Install Redis Extension in Docker Desktop
3. Start redis by running `docker run --name my-redis -p 6379:6379 -d redis` on your terminal

### Get Started 
1. Clone the repository: `git clone https://github.com/yamaceay/rediscache`
2. Change the current directory: `cd rediscache`
3. Make sure the Redis container is running on Docker Desktop (-> Containers)
4. Start the server: `go run .` (optionally you can add parameters: `--dbAddress=<dbHost>:<dbPort> --ipAddress=<ipHost>:<ipPort> --ttlMinutes:<ttlMinutes>`
5. Starting a new client: 
  1. Open a new terminal in the same directory
  2. Run `go run . --mode=client` to see the commands

### CLI Arguments
#### Server Mode

| Name | Description                     | Type | Options | Default |
| :----: | :-----------------------------: | :-: | :-------: | :-------: |
| dbAddress | Database address | string | any ip address | "localhost:6379" |
| ipAddress | Application address | string | any ip address | "localhost:8080" |
| ttlMinutes | Expiry time of data objects in minutes | int | any positive integer | 10080 |    

#### Client Mode

| Name | Alias | Description                     | Type | Options | Default |
| :----: | :---: | :-----------------------------: | :-: | :-------: | :-------: |
| method | X | What to do | string | "", "get", "set" | "" |
| key, value | k, v | Key-value pairs | string | - | "" |
| db | - | Which database | int | 0, 1, … | 0 |