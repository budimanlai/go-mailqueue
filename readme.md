# How To Use
**1. Create config folder**

**2. Create main.conf file in config folder**
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
// create dist folder

// build
./scripts/compile.sh
```

**4. Run**
```
// run service and show the output in terminal
./dist/mail-queue-{GIT_HEAD_CODE} run

// run service in background
./dist/mail-queue-{GIT_HEAD_CODE} start

// stop the service
./dist/mail-queue-{GIT_HEAD_CODE} stop

// help
./dist/mail-queue-{GIT_HEAD_CODE} help
```