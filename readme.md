# How To Use
**1. Create config folder**

**2. Create config file in config/main.conf**
```
[database]
hostname=127.0.0.1
username=root
password=
database=your_database
port=3306

[smtp]
hostname=smtp_hostname
port=smtp_port
username=smtp_username
password=smtp_password
tls=true
ping=12
idle=30
```

**3. Import mail_queue.sql to database**
**4. Build (for linux)**
```
./scripts/compile.sh
```

**4. Run**
```
// run service and show the output in terminal
./mail-queue run

// run service in background
./mail-queue start

// stop the service
./mail-queue stop

// help
./mail-queue help
```