## LostBot - A game in text

This project is designed to test my knowledge of the GO language, using a database on PostgresSQL and Docker

You can use it as a sample project, and you can also make pull requests to add new functionality

For details of changes in the bot go to [CHANGELOG.md](https://github.com/AndreyFulov/lostbot-go/blob/main/CHANGELOG.md)

To run this you need to

```cmd
docker build -t lostbot .
```

and

```cmd
docker-compose up
```

Now bot will work!

### Road Map

- [x] First version of players table
- [x] System of Casino
- [x] Top money
- [x] Display profile command
- [ ] Business System
  - [x] AFK Earning
  - [x] Database for business
  - [ ] Upgrade system
  - [x] Buy business
- [ ] Shares system
  - [ ] Shares can fall and rise thanks to player investment
  - [ ] Players can add their own shares on market
