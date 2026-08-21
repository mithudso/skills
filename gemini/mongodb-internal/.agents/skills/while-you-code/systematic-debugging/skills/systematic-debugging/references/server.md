# MongoDB Server Debugging Context

When debugging in `10gen/mongo` or `10gen/mms`:

- Run tests with resmoke: `python buildscripts/resmoke.py run --suite <suite> <test-file>`
- Check recent Evergreen failures: `evergreen patch -p mongodb-mongo-master`
- Core dump analysis: `gdb --args mongod <coredump>` then `bt full`
- Log location: `data/db/mongod.log` (local) or fetch from Evergreen task artifacts
- Common failure patterns: WT cache pressure, replication lag, index build failures
