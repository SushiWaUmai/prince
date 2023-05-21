# Prince
A simple whatsapp bot written in golang using whatsmeow.


## How to use

Sample docker-compose.yml file
```
version: "3"
services:
  prince:
    image: sushiwaumai/prince
    container_name: prince
    volumes:
      - ./prince.db:/prince.db
    env_file: .env
```

## License
This project is licensed under the [MIT](./LICENSE) license
